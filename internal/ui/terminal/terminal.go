package terminal

import (
	"fmt"
	"strings"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vimgo/vimgo/internal/board"
	"github.com/vimgo/vimgo/internal/game"
	"github.com/vimgo/vimgo/internal/rules"
	"github.com/vimgo/vimgo/internal/sgf"
	"github.com/vimgo/vimgo/internal/vim"
)

var (
	boardStyle = lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63"))

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("63")).
			Padding(0, 1)

	blackStone = "●"
	whiteStone = "○"
	// Box drawing characters
	topLeft     = "┌"
	topRight    = "┐"
	bottomLeft  = "└"
	bottomRight = "┘"
	horizontal  = "─"
	vertical    = "│"
	intersection = "┼"
	teeLeft     = "├"
	teeRight    = "┤"
	teeTop      = "┬"
	teeBottom   = "┴"
	starPoint   = "╋"
)

var (
	gridColor = lipgloss.Color("240") // Subtle grey/blue
	starColor = lipgloss.Color("214") // Orange/Yellow
)

type Model struct {
	Game    *game.Game
	Handler *vim.Handler
	Error   error
	Width      int
	Height     int
	ShowCoords bool
	ShowHelp   bool
	ScoreText  string
}

func NewModel(size int) Model {
	return Model{
		Game:    game.NewGame(size),
		Handler: vim.NewHandler(size),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		key := msg.String()
		
		// Map bubbletea keys to our handler strings
		switch key {
		case "up": key = "k"
		case "down": key = "j"
		case "left": key = "h"
		case "right": key = "l"
		case "backspace": key = "backspace"
		case "esc": key = "esc"
		case "enter": key = "enter"
		}

		action := m.Handler.HandleKey(key)
		if action != nil {
			switch action.Type {
			case vim.ActionPlaceStone:
				m.Error = m.Game.Move(m.Handler.CursorX, m.Handler.CursorY)
			case vim.ActionUndo:
				m.Error = m.Game.Undo()
			case vim.ActionCommand:
				return m.handleCommand(action.Value)
			}
		}

		if key == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) handleCommand(cmd string) (tea.Model, tea.Cmd) {
	m.Error = nil
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return m, nil
	}

	switch parts[0] {
	case "q", "quit":
		if m.ShowHelp {
			m.ShowHelp = false
			return m, nil
		}
		return m, tea.Quit
	case "e", "edit":
		filename := "game.sgf"
		if len(parts) > 1 {
			filename = parts[1]
		}
		err := m.loadSGF(filename)
		if err != nil {
			m.Error = err
		}
	case "w", "write":
		filename := "game.sgf"
		if len(parts) > 1 {
			filename = parts[1]
		}
		err := m.saveSGF(filename)
		if err != nil {
			m.Error = err
		}
	case "undo":
		m.Error = m.Game.Undo()
	case "pass":
		m.Game.CurrentPlayer = m.Game.CurrentPlayer.Opposite()
	case "coords", "coordinates", "c":
		m.ShowCoords = !m.ShowCoords
	case "?", "help":
		m.ShowHelp = true
	case "score":
		method := "chinese"
		if len(parts) > 1 {
			method = parts[1]
		}
		// Basic komi 7.5 for Chinese, 6.5 for Japanese? Let's assume 7.5 for now or 0.
		// The user didn't specify komi, let's use 7.5
		score := rules.CountScore(m.Game.Board, method, m.Game.BlackCaptures, m.Game.WhiteCaptures, 7.5)
		m.ScoreText = fmt.Sprintf("[W %.1f B %.1f]", score.White, score.Black)
	default:
		m.Error = fmt.Errorf("unknown command: %s", parts[0])
	}
	return m, nil
}

func (m Model) saveSGF(filename string) error {
	content := sgf.SimpleSGFWriter(m.Game.Board.Size, m.Game.Moves)
	return os.WriteFile(filename, []byte(content), 0644)
}

func (m *Model) loadSGF(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	size, moves, err := sgf.ParseSGF(string(content))
	if err != nil {
		return err
	}

	// Initialize new game
	newGame := game.NewGame(size)
	
	// Replay moves
	for _, moveStr := range moves {
		// Parse move string like "B[pd]"
		if len(moveStr) < 4 || moveStr[1] != '[' || moveStr[len(moveStr)-1] != ']' {
			continue // Skip malformed
		}
		color := string(moveStr[0])
		val := moveStr[2 : len(moveStr)-1]
		
		// Handle pass
		if val == "" {
			if color == "B" {
				newGame.CurrentPlayer = board.Black
			} else {
				newGame.CurrentPlayer = board.White
			}
			newGame.CurrentPlayer = newGame.CurrentPlayer.Opposite()
			continue
		}

		x, y, err := sgf.FromSGFCoord(val)
		if err != nil {
			return err
		}

		// Set correct player for move
		if color == "B" {
			newGame.CurrentPlayer = board.Black
		} else {
			newGame.CurrentPlayer = board.White
		}

		if err := newGame.Move(x, y); err != nil {
			return fmt.Errorf("failed to replay move %s: %v", moveStr, err)
		}
	}

	m.Game = newGame
	m.Handler = vim.NewHandler(size)
	return nil
}

func (m Model) View() string {
	if m.Width == 0 {
		return "Initializing..."
	}

	var s strings.Builder

	// Header
	header := lipgloss.NewStyle().Bold(true).Render("VimGo - Go with Vim keybindings")
	s.WriteString(lipgloss.NewStyle().Width(m.Width).Align(lipgloss.Center).Render(header))
	s.WriteString("\n")

	// Render board content
	var boardView strings.Builder
	size := m.Game.Board.Size
	// Helper for coordinates
	getCoordChar := func(i int) string {
		if i >= 8 { i++ } // Skip 'I'
		return string(rune('A' + i))
	}
	
	coordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// Top Coordinates
	if m.ShowCoords {
		boardView.WriteString("   ") // Space for left numbers
		for x := 0; x < size; x++ {
			boardView.WriteString(coordStyle.Render(getCoordChar(x)))
			if x < size-1 {
				boardView.WriteString("   ")
			}
		}
		boardView.WriteString("\n")
	}

	for y := 0; y < size; y++ {
		// Left Coordinate
		if m.ShowCoords {
			num := fmt.Sprintf("%2d ", size-y)
			boardView.WriteString(coordStyle.Render(num))
		}

		for x := 0; x < size; x++ {
			// Determine grid character based on position
			char := intersection
			if y == 0 {
				if x == 0 { char = topLeft
				} else if x == size-1 { char = topRight
				} else { char = teeTop }
			} else if y == size-1 {
				if x == 0 { char = bottomLeft
				} else if x == size-1 { char = bottomRight
				} else { char = teeBottom }
			} else {
				if x == 0 { char = teeLeft
				} else if x == size-1 { char = teeRight }
			}

			// Star points (Hoshi)
			isStar := false
			if size == 19 {
				isStar = (x == 3 || x == 9 || x == 15) && (y == 3 || y == 9 || y == 15)
			}
			if size == 13 {
				isStar = (x == 3 || x == 6 || x == 9) && (y == 3 || y == 6 || y == 9)
			}
			if size == 9 {
				isStar = (x == 2 || x == 6) && (y == 2 || y == 6) || (x == 4 && y == 4)
			}
			if isStar {
				char = starPoint
			}

			c := m.Game.Board.At(x, y)
			cellContent := ""
			if c == board.Black {
				cellContent = blackStone
			} else if c == board.White {
				cellContent = whiteStone
			} else {
				// Render grid character with color
				style := lipgloss.NewStyle().Foreground(gridColor)
				if isStar {
					style = style.Foreground(starColor)
				}
				cellContent = style.Render(char)
			}

			// Cursor logic
			if m.Handler.CursorX == x && m.Handler.CursorY == y {
				style := lipgloss.NewStyle().Background(lipgloss.Color("201")).Foreground(lipgloss.Color("15"))
				boardView.WriteString(style.Render(cellContent))
			} else {
				boardView.WriteString(cellContent)
			}

			// Horizontal connection line
			if x < size-1 {
				lineChar := horizontal + horizontal + horizontal
				style := lipgloss.NewStyle().Foreground(gridColor)
				boardView.WriteString(style.Render(lineChar))
			}
		}
		
		// Right Coordinate
		if m.ShowCoords {
			num := fmt.Sprintf(" %-2d", size-y)
			boardView.WriteString(coordStyle.Render(num))
		}

		boardView.WriteString("\n")

		// Vertical connection row (for square look)
		if y < size-1 {
			// Left spacer for coords
			if m.ShowCoords {
				boardView.WriteString("   ")
			}

			for x := 0; x < size; x++ {
				style := lipgloss.NewStyle().Foreground(gridColor)
				boardView.WriteString(style.Render(vertical))
				if x < size-1 {
					boardView.WriteString("   ") // Spacing between vertical lines
				}
			}
			boardView.WriteString("\n")
		}
	}

	// Bottom Coordinates
	if m.ShowCoords {
		boardView.WriteString("   ") // Space for left numbers
		for x := 0; x < size; x++ {
			boardView.WriteString(coordStyle.Render(getCoordChar(x)))
			if x < size-1 {
				boardView.WriteString("   ")
			}
		}
		boardView.WriteString("\n")
	}

	// Center the board
	styledBoard := boardStyle.Render(boardView.String())

	if m.ShowHelp {
		helpText := "\n  VimGo Help\n\n"
		helpText += "  hjkl    Move cursor\n"
		helpText += "  x       Place stone\n"
		helpText += "  u       Undo\n"
		helpText += "  i       Insert Mode\n"
		helpText += "  :w      Save (game.sgf)\n"
		helpText += "  :c      Toggle Coords\n"
		helpText += "  :e [f]  Load SGF\n"
		helpText += "  :score  [chinese|japanese]\n"
		helpText += "  :?      Show Help\n"
		helpText += "  :q      Quit / Close Help\n"
		
		helpBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(helpText)
		
		// Overlay help box or replace board? Let's replace for now as overlay is tricky with lipgloss text only
		// Actually lipgloss.Place effectively centers, so we can just swap what we center.
		styledBoard = helpBox
	}

	centeredBoard := lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, styledBoard)
	s.WriteString(centeredBoard)

	// Error message area
	if m.Error != nil {
		errorText := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Error: %v", m.Error))
		s.WriteString("\n" + errorText)
	}

	// Status bar at the bottom
	modeStr := m.Handler.Mode.String()
	coord := game.CoordinateToString(m.Game.Board.Size, m.Handler.CursorX, m.Handler.CursorY)
	turn := "Black"
	if m.Game.CurrentPlayer == board.White {
		turn = "White"
	}
	
	statusText := fmt.Sprintf(" -- %s -- %dx%d -- %s -- Turn: %d -- [%s] -- %s", 
		modeStr, size, size, turn, len(m.Game.History)+1, coord, m.ScoreText)
	
	if m.Handler.Mode == vim.Command {
		statusText = ":" + m.Handler.CommandBuffer
	}

	// Pin to bottom
	statusBar := statusBarStyle.Width(m.Width).Render(statusText)
	return lipgloss.JoinVertical(lipgloss.Top, s.String(), statusBar)
}
