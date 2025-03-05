package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetCountryInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	// Extract the country code from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		respondWithError(w, http.StatusBadRequest, "Country code is missing in the request")
		return
	}
	countryCode := strings.TrimSpace(strings.ToLower(pathParts[4]))

	// Fetch country data from RestCountriesAPI
	urlAPI := RestCountriesAPI + "alpha/" + countryCode
	resp, err := http.Get(urlAPI)
	if err != nil {
		respondWithError(w, http.StatusServiceUnavailable, "Failed to fetch country data")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading response from API")
		return
	}

	// Unmarshal into a slice
	var countryData []RestCountriesData
	if err := json.Unmarshal(body, &countryData); err != nil {
		fmt.Println("JSON Unmarshal Error:", err)
		respondWithError(w, http.StatusInternalServerError, "Error processing country data")
		return
	}

	// Ensure at least one country exists in the response
	if len(countryData) == 0 {
		respondWithError(w, http.StatusNotFound, "No country data found")
		return
	}

	// Use the first country in the response
	selectedCountry := countryData[0]

	// Fetch cities from CountriesNowAPI
	citiesAPI := fmt.Sprintf("%s/countries/cities", CountriesNowAPI)
	citiesRequestBody := fmt.Sprintf(`{"country": "%s"}`, selectedCountry.Name.CountryName)

	citiesResp, err := http.Post(citiesAPI, "application/json", strings.NewReader(citiesRequestBody))
	if err != nil {
		fmt.Println("Failed to fetch cities:", err)
		respondWithError(w, http.StatusServiceUnavailable, "Failed to fetch cities data")
		return
	}
	defer citiesResp.Body.Close()

	citiesBody, err := ioutil.ReadAll(citiesResp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading cities response")
		return
	}

	// Parse Cities API response
	var citiesResponse struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(citiesBody, &citiesResponse); err != nil {
		fmt.Println("JSON Unmarshal Error (Cities):", err)
		respondWithError(w, http.StatusInternalServerError, "Error processing cities data")
		return
	}

	// Construct final JSON response with full data
	response := map[string]interface{}{
		"iso_code":   strings.ToUpper(countryCode),
		"name":       selectedCountry.Name.CountryName,
		"official":   selectedCountry.Name.Official,
		"continents": selectedCountry.Continents,
		"population": selectedCountry.Population,
		"languages":  selectedCountry.Languages,
		"bordering":  selectedCountry.Bordering,
		"flag":       selectedCountry.CountryFlag.Png,
		"capital":    "",
		"currencies": selectedCountry.Currencies,
		"region":     selectedCountry.Region,
		"subregion":  selectedCountry.Subregion,
		"timezones":  selectedCountry.Timezones,
		"cities":     citiesResponse.Data,
	}

	if len(selectedCountry.Capital) > 0 {
		response["capital"] = selectedCountry.Capital[0]
	}

	// Convert response to JSON
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating JSON response")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// respondWithError sends an error response with a specific HTTP status code
func respondWithError(w http.ResponseWriter, errorCode int, message string) {
	errorResponse := ErrorResponse{
		ErrorCode: errorCode,
		Message:   message,
	}
	responseJSON, _ := json.Marshal(errorResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	w.Write(responseJSON)
}
