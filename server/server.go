package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Chat struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newChat() *Chat {
	return &Chat{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (c *Chat) run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)
			}
		case message := <-c.broadcast:
			log.Printf("received message: %s", message)
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}
		}
	}
}

const (
	// maxMessageSize is the maximum size in bytes of a message the server
	maxMessageSize = 512

	// writeWait is the maximum time in seconds to write a message to the peer.
	writeWait = 10 * time.Second

	// pongWait is the maximum time in seconds to wait for a ping from the
	// client before assuming the connection is dead.
	pongWait = 60 * time.Second

	// pingPeriod - send ping every 54 seconds to client
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	name string
	chat *Chat
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.chat.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.chat.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			msg := fmt.Sprintf("%s: %s", c.name, message)

			w.Write([]byte(msg))

			log.Printf("received message: %s", msg)
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Chat) connectUser(u User, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if u.Name == "" {
		http.Error(w, fmt.Sprintf("Name is required, cannot find name of user with ID %d", u.ID), http.StatusBadRequest)
		return
	}

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{name: u.Name, chat: c, conn: conn, send: make(chan []byte, 256)}
	log.Printf("%s connected: %s", client.name, client.conn.RemoteAddr())
	client.chat.register <- client

	go client.writePump()
	go client.readPump()
}
