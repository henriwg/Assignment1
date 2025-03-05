package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Start time for calculating uptime
var startTime time.Time

func init() {
	startTime = time.Now()
}

// StatusHandler handles the /status endpoint request
func GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check status for both APIs
	countriesNowStatus := getAPIStatus(CountriesNowAPI + "/countries/")
	restCountriesStatus := getAPIStatus(RestCountriesAPI + "all")

	// Prepare response
	statusInfo := StatusResponse{
		CountriesNowAPI:  fmt.Sprintf("%d", countriesNowStatus),
		RestCountriesAPI: fmt.Sprintf("%d", restCountriesStatus),
		Version:          "v1",
		Uptime:           int(time.Since(startTime).Seconds()),
	}

	// Convert response to JSON with indentation
	responseJSON, err := json.MarshalIndent(statusInfo, "", "  ")
	if err != nil {
		log.Printf("Encoding error for statusInfo: %v", err)
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// getAPIStatus makes a request to an external API and returns the HTTP status code
func getAPIStatus(url string) int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v", url, err)
		return 503
	}

	req.Header.Add("Authorization", "Bearer YOUR_API_KEY")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching API status from %s: %v", url, err)
		return 503
	}
	defer resp.Body.Close()

	log.Printf("Checked API status: %s -> %d", url, resp.StatusCode)
	return resp.StatusCode
}
