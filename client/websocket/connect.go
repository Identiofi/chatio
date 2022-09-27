package websocket

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type session struct {
	ws      *websocket.Conn
	rl      *readline.Instance
	errChan chan error
}

func bytesToFormattedHex(bytes []byte) string {
	text := hex.EncodeToString(bytes)
	return regexp.MustCompile("(..)").ReplaceAllString(text, "$1 ")
}

func (s *session) read() {
	rxSprintf := color.New(color.FgGreen).SprintfFunc()

	for {
		t, msg, err := s.ws.ReadMessage()
		if err != nil {
			s.errChan <- err
			return
		}

		var text string
		switch t {
		case websocket.TextMessage:
			text = string(msg)
		case websocket.BinaryMessage:
			bytesToFormattedHex(msg)
		default:
			s.errChan <- fmt.Errorf("unknown websocket frame type: %d", t)
		}

		fmt.Fprintf(s.rl.Stdout(), rxSprintf("< %s\n", text))
	}
}

func (s *session) write() {
	for {
		line, err := s.rl.Readline()
		if err != nil {
			s.errChan <- err
			return
		}

		err = s.ws.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			s.errChan <- err
			return
		}
	}
}

func Connect(url string, rlConf *readline.Config) error {
	headers := make(http.Header)
	headers.Add("Origin", url)

	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	ws, _, err := dialer.Dial(url, headers)
	if err != nil {
		return err
	}
	rl, err := readline.NewEx(rlConf)
	if err != nil {
		return err
	}
	defer rl.Close()

	s := &session{
		ws:      ws,
		rl:      rl,
		errChan: make(chan error),
	}

	go s.read()
	go s.write()
	return <-s.errChan
}
