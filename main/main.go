package main

import (
	"CloudTechnologiesHenrikWG/server"
	"log"
	"net/http"
)

func main() {
	mux := server.SetupRouter()
	log.Println("Server is starting on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
