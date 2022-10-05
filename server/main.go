package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

//go:embed static
var staticFiles embed.FS

func main() {
	chat := newChat()
	go chat.run()

	staticFS := http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	// Serve static files
	http.Handle("/static/", fs)

	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/register", registerUserHandler)
	http.HandleFunc("/unregister", unregisterUserHandler)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(chat, w, r)
	})

	log.Printf("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

type Message struct {
	Message string `json:"message"`
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

var userList = []User{
	{Name: "John Doe", ID: 1},
	{Name: "Jane Doe", ID: 2},
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	noID := make([]User, len(userList))
	for i, u := range userList {
		noID[i] = User{Name: u.Name, ID: 0}
	}
	json.NewEncoder(w).Encode(noID)
}

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	json.NewDecoder(r.Body).Decode(&newUser)

	if r.Method != "POST" || newUser.Name == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// generate random ID, add length of userList to ensure uniqueness
	newUser.ID = rand.Intn(1000000) + len(userList)

	// append user to userList
	userList = append(userList, newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{ "id": "%d" }`, newUser.ID)
}

func unregisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var usr User
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

func serveWebsocket(chat *Chat, w http.ResponseWriter, r *http.Request) {
	in := r.URL.Query().Get("id")

	var id int
	id, err := strconv.Atoi(in)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check if user is registered
	for _, u := range userList {
		if u.ID == id {
			log.Printf("Connecting user %s ...", u.Name)
			chat.connectUser(u, w, r)
			return
		}
	}

	http.Error(w, `{ "message": "user not found"}`, http.StatusNotFound)
}
