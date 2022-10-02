package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	// c, err := newClient(URL)

	// c.Run()
}

type client struct {
	// the websocket connection
	conn *websocket.Conn
	// errChan is used to send errors to the main goroutine
	// from the listen goroutine
	errChan chan error
}

func (c *client) write() {
	// your code here
}

func (c *client) listen() {
	// your code here
}

// run starts the client
func (c *client) run() {
	// interrupt may be sent by the OS when the user presses Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// start listen and write in separate goroutines
	// allows us to read and write concurrently
	go c.listen()
	go c.write()

	// wait for an error or interrupt, this is done on the main thread
	// so that we can block until the program is terminated
	for {
		select {
		case err := <-c.errChan:
			c.Close()
			fail("error: %v\n", err)
		case <-interrupt:
			c.Close()
			return
		}
	}
}

func (c *client) Close() {
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func fail(msg string, o ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, o...)
	os.Exit(1)
}
