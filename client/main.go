package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"gopkg.in/alecthomas/kingpin.v2"
)

var url = kingpin.Flag("url", "websocket url").Required().String()

func main() {
	kingpin.Parse()

	c, err := newClient(*url)
	if err != nil {
		fail("error: %v", err)
	}
	c.run()
}

type client struct {
	conn    *websocket.Conn
	errChan chan error
}

func newClient(url string) (*client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return &client{
		conn:    conn,
		errChan: make(chan error),
	}, nil
}

func (c *client) write() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		c.conn.WriteMessage(websocket.TextMessage, []byte(stdin.Text()))
	}
}

func (c *client) listen() {
	// a infinite loop to read messages from the server
	for {
		t, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.errChan <- err
			return
		}

		var text string
		switch t {
		case websocket.TextMessage:
			text = string(msg)
		default:
			c.errChan <- fmt.Errorf("unknown websocket frame type: %d", t)
		}
		fmt.Fprintf(os.Stdout, "%s\n", text)
	}
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
