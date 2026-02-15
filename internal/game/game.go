package game

import (
	"fmt"
	"github.com/vimgo/vimgo/internal/board"
	"github.com/vimgo/vimgo/internal/rules"
)

type Game struct {
	Board         *board.Board
	CurrentPlayer board.Color
	History       []*board.Board
	undoStack     []undoState
	BlackCaptures int
	WhiteCaptures int
	LastMove      *board.Point
	Moves         []string // Store moves in SGF format: B[pd], W[aa]
}

type undoState struct {
	currentPlayer board.Color
	blackCaptures int
	whiteCaptures int
	lastMove      *board.Point
	movesLen      int
}

func NewGame(size int) *Game {
	return &Game{
		Board:         board.New(size),
		CurrentPlayer: board.Black,
		History:       []*board.Board{},
		undoStack:     []undoState{},
	}
}

// Move places a stone at (x, y) if valid, updates captures and turn.
func (g *Game) Move(x, y int) error {
	if !rules.IsMoveValid(g.Board, x, y, g.CurrentPlayer) {
		return fmt.Errorf("invalid move at (%d, %d)", x, y)
	}

	// Ko detection (basic version: compare to immediate previous state)
	// For full Ko, we'd check against all previous states or use Zobrist hashing.
	// We check if the resulting board would be identical to the one before the previous move.
	tempBoard := g.Board.Copy()
	tempBoard.Set(x, y, g.CurrentPlayer)
	captured := rules.FindCapturedStones(tempBoard, x, y, g.CurrentPlayer)
	for _, p := range captured {
		tempBoard.Set(p.X, p.Y, board.Empty)
	}

	if len(g.History) > 0 {
		prevBoard := g.History[len(g.History)-1]
		if boardsEqual(tempBoard, prevBoard) {
			return fmt.Errorf("ko violation")
		}
	}

	// Apply move
	var prevLastMove *board.Point
	if g.LastMove != nil {
		prevLastMove = &board.Point{X: g.LastMove.X, Y: g.LastMove.Y}
	}
	g.History = append(g.History, g.Board.Copy())
	g.undoStack = append(g.undoStack, undoState{
		currentPlayer: g.CurrentPlayer,
		blackCaptures: g.BlackCaptures,
		whiteCaptures: g.WhiteCaptures,
		lastMove:      prevLastMove,
		movesLen:      len(g.Moves),
	})
	g.Board.Set(x, y, g.CurrentPlayer)

	// Update captures
	if g.CurrentPlayer == board.Black {
		g.BlackCaptures += len(captured)
	} else {
		g.WhiteCaptures += len(captured)
	}

	// Remove captured stones
	for _, p := range captured {
		g.Board.Set(p.X, p.Y, board.Empty)
	}

	// Update last move and switch player
	g.LastMove = &board.Point{X: x, Y: y}

	// Record move for SGF
	colorStr := "B"
	if g.CurrentPlayer == board.White {
		colorStr = "W"
	}
	// Note: this assumes we have sgf package or just format it here.
	// We'll format it here for simplicity since game doesn't strictly depend on sgf for recording.
	sgfX := rune('a' + x)
	sgfY := rune('a' + y)
	g.Moves = append(g.Moves, fmt.Sprintf("%s[%c%c]", colorStr, sgfX, sgfY))

	g.CurrentPlayer = g.CurrentPlayer.Opposite()

	return nil
}

// Undo reverts to the previous board state.
func (g *Game) Undo() error {
	if len(g.History) == 0 {
		return fmt.Errorf("nothing to undo")
	}
	if len(g.undoStack) == 0 {
		return fmt.Errorf("undo state corrupted")
	}
	prev := g.undoStack[len(g.undoStack)-1]
	if prev.movesLen < 0 || prev.movesLen > len(g.Moves) {
		return fmt.Errorf("undo state corrupted")
	}

	g.Board = g.History[len(g.History)-1]
	g.History = g.History[:len(g.History)-1]
	g.undoStack = g.undoStack[:len(g.undoStack)-1]
	g.CurrentPlayer = prev.currentPlayer
	g.BlackCaptures = prev.blackCaptures
	g.WhiteCaptures = prev.whiteCaptures
	g.LastMove = prev.lastMove
	g.Moves = g.Moves[:prev.movesLen]
	return nil
}

func boardsEqual(b1, b2 *board.Board) bool {
	if b1.Size != b2.Size {
		return false
	}
	for i := range b1.Grid {
		if b1.Grid[i] != b2.Grid[i] {
			return false
		}
	}
	return true
}

func CoordinateToString(size, x, y int) string {
	// Column: A, B, C... (skip I)
	col := rune('A' + x)
	if x >= 8 { // Skip 'I'
		col++
	}
	// Row: 1 is bottom
	row := size - y
	return fmt.Sprintf("%c%d", col, row)
}
