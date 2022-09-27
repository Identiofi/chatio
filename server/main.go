package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	// http.HandleFunc("/users", usersHandler)

	log.Printf("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

type message struct {
	Message string `json:"message"`
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// init new variable of type message
	var msg message
	// decode the request body into the message variable,
	// passing a pointer to the message variable
	json.NewDecoder(r.Body).Decode(&msg)

	if msg.Message != "Hello world" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// encode the message variable as JSON and write it to the response writer.
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message{Message: "Hello World!"})
}

// TODO: type user struct

// var userList = []user{}
