package main

import (
	"log"
	"net/http"
)

func main() {
	chat := newChat()
	go chat.run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWs(chat, w, r)
	})

	log.Printf("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
