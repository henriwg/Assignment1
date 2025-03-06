package main

import (
	"CloudTechnologiesHenrikWG/handlers"
	"log"
	"net/http"
)

func main() {
	mux := SetupRouter()
	log.Println("Server is starting on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// Register routes (prefix patterns for paths with parameters)
	mux.HandleFunc("/countryinfo/v1/info/", handlers.GetCountryInfo)
	mux.HandleFunc("/countryinfo/v1/population/", handlers.GetPopulation)
	mux.HandleFunc("/countryinfo/v1/status", handlers.GetStatus)
	return mux
}
