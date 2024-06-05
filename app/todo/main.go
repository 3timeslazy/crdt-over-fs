package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jessevdk/go-flags"
)

type AppOptions struct {
	Device string `short:"d" long:"device" required:"true" description:"Device ID"`
}

func main() {
	opts := AppOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	prog := tea.NewProgram(NewApp(opts.Device))
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
