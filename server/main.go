package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

//go:embed static
var embeddedFiles embed.FS

func main() {
	chat := newChat()
	go chat.run()

	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/register", registerUserHandler)
	http.HandleFunc("/unregister", unregisterUserHandler)
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWebsocket(chat, w, r)
	})

	log.Printf("Starting server on port 8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(readFile("./static/index.html")))
}

type message struct {
	Message string `json:"message"`
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

	id, err := strconv.Atoi(in)
	if err != nil {
		log.Fatalf("could not convert id to int: %v", err)
	}

	// check if user is registered
	for _, u := range userList {
		if u.ID == id {
			chat.connectUser(u, w, r)
			return
		}
	}

	http.Error(w, `{ "message": "user not found"}`, http.StatusNotFound)
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return content
}
