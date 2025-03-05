package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetPopulation retrieves population data for a country using its ISO code
func GetPopulation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Extract the country ISO2 code from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		respondWithError(w, http.StatusBadRequest, "Country ISO code is missing in the request")
		return
	}
	iso2Code := strings.TrimSpace(strings.ToUpper(pathParts[4]))

	// Extract the limit query parameter (if provided)
	limitStr := r.URL.Query().Get("limit")
	startYear, endYear, limitErr := parseLimit(limitStr)

	restCountriesURL := fmt.Sprintf("%salpha/%s", RestCountriesAPI, iso2Code)

	resp, err := http.Get(restCountriesURL)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Failed to fetch country name")
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading response from API")
		return
	}

	// Parse RestCountries API response
	var countryDataArray []RestCountriesData
	if err := json.Unmarshal(body, &countryDataArray); err != nil {
		log.Println("JSON Unmarshal Error (RestCountries):", err)
		respondWithError(w, http.StatusInternalServerError, "Error processing country data")
		return
	}

	// Ensure at least one country exists in the response
	if len(countryDataArray) == 0 {
		respondWithError(w, http.StatusNotFound, "No country data found")
		return
	}

	// Use the first country in the response
	selectedCountry := countryDataArray[0]
	countryName := selectedCountry.Name.CountryName
	if countryName == "" {
		respondWithError(w, http.StatusNotFound, "Could not find country name for the provided ISO code")
		return
	}

	// ----------------------------------------- CountriesNow API ---------------------------------------------------
	// Use the retrieved country name to fetch population data
	populationAPI := fmt.Sprintf("%s/countries/population", CountriesNowAPI)
	payload := map[string]string{"country": countryName}
	jsonPayload, _ := json.Marshal(payload)

	// Send the request
	popResp, err := http.Post(populationAPI, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Failed to fetch population data")
		return
	}
	defer popResp.Body.Close()

	// Read population API response
	popBody, err := ioutil.ReadAll(popResp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading population API response")
		return
	}

	// Unmarshal response
	var popResponse struct {
		Error bool   `json:"error"`
		Msg   string `json:"msg"`
		Data  struct {
			PopulationCounts []struct {
				Year  int `json:"year"`
				Value int `json:"value"`
			} `json:"populationCounts"`
		} `json:"data"`
	}
	if err := json.Unmarshal(popBody, &popResponse); err != nil {
		log.Println("JSON Unmarshal Error (Population):", err)
		respondWithError(w, http.StatusInternalServerError, "Error processing population data")
		return
	}

	// If API returns an error
	if popResponse.Error {
		respondWithError(w, http.StatusNotFound, "Population data not found")
		return
	}

	// Filter population data based on the limit parameter
	filteredPopulation := []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	}{}

	var sum, count int

	for _, pop := range popResponse.Data.PopulationCounts {
		if !limitErr {
			if pop.Year < startYear || pop.Year > endYear {
				continue
			}
		}
		filteredPopulation = append(filteredPopulation, struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		}{Year: pop.Year, Value: pop.Value})
		sum += pop.Value
		count++
	}

	meanPopulation := 0
	if count > 0 {
		meanPopulation = sum / count
	}

	// Construct response JSON
	populationData := PopulationResponse{
		Mean:   meanPopulation,
		Values: filteredPopulation,
	}

	// Convert response to JSON
	responseJSON, err := json.MarshalIndent(populationData, "", "  ")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating JSON response")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// parseLimit extracts the start and end years from the limit parameter
func parseLimit(limitStr string) (int, int, bool) {
	if limitStr == "" {
		return 0, 0, true
	}

	parts := strings.Split(limitStr, "-")
	if len(parts) != 2 {
		return 0, 0, true
	}

	startYear, err1 := strconv.Atoi(parts[0])
	endYear, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil || startYear > endYear {
		return 0, 0, true
	}

	return startYear, endYear, false
}
