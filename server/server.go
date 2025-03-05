package server

import (
	"CloudTechnologiesHenrikWG/handlers"
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// Register routes (prefix patterns for paths with parameters)
	mux.HandleFunc("/countryinfo/v1/info/", handlers.GetCountryInfo)
	mux.HandleFunc("/countryinfo/v1/population/", handlers.GetPopulation)
	mux.HandleFunc("/countryinfo/v1/status", handlers.GetStatus)
	return mux
}
