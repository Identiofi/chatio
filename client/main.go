package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/identiofi/chatio/websocket"
	"github.com/spf13/cobra"
)

func root(cmd *cobra.Command, args []string) {
	dest := "ws://localhost:8080/chat"

	err := websocket.Connect(dest, &readline.Config{
		Prompt: "> ",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err != io.EOF && err != readline.ErrInterrupt {
			os.Exit(1)
		}
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "chat URL",
		Short: "CLI chat tool",
		Run:   root,
	}

	rootCmd.Execute()
}
