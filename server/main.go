package main

import (
	"encoding/json"
	"fmt"
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

type Message struct {
	Message string `json:"message"`
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// init new variable of type message
	var msg Message
	// decode the request body into the message variable,
	// passing a pointer to the message variable
	json.NewDecoder(r.Body).Decode(&msg)

	// check if the message equals "hello World"
	if msg.Message != "Hello World" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// return a response with the message "Message received"
	// return the content type as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{ "message": "Message received"}`)
}

// TODO - TASK 1: type user struct

// userList is a variable that stores all users of type user
// as it's outside of the function scope, it's accessible to all functions (global)
// Uncomment the line below to use it
// var userList = []User{{ Name: "John Doe", ID: 1234 }, { Name: "Jane Doe", ID: 5678 }}

// TODO: implement usersHandler
