package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vimgo/vimgo/internal/ui/terminal"
)

func main() {
	size := flag.Int("size", 19, "Board size (9, 13, or 19)")
	flag.Parse()

	if *size != 9 && *size != 13 && *size != 19 {
		fmt.Println("Invalid size. Please use 9, 13, or 19.")
		os.Exit(1)
	}

	m := terminal.NewModel(*size)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running VimGo: %v", err)
		os.Exit(1)
	}
}
