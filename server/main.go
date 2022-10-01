package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/unregister", unregisterUser)

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
type user struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

var userList = []user{
	{Name: "John Doe", ID: "1"},
	{Name: "Jane Doe", ID: "2"},
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userList)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var newUser user
	json.NewDecoder(r.Body).Decode(&newUser)

	if r.Method != "POST" || newUser.Name == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// generate random ID, add length of userList to ensure uniqueness
	newUser.ID = strconv.Itoa(rand.Intn(1000000) + len(userList))

	// append user to userList
	userList = append(userList, newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{ "id": "%s" }`, newUser.ID)
}

func unregisterUser(w http.ResponseWriter, r *http.Request) {
	var usr user
	json.NewDecoder(r.Body).Decode(&usr)

	for i, u := range userList {
		if u.ID == usr.ID {
			userList = append(userList[:i], userList[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{ "message": "user unreigistered" }`)
			return
		}
	}

	http.Error(w, `{ "message": "user not found"}`, http.StatusNotFound)
}
