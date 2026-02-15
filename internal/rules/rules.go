package rules

import (
	"github.com/vimgo/vimgo/internal/board"
)

// Group represents a set of connected stones of the same color.
type Group struct {
	Stones    []board.Point
	Liberties []board.Point
	Color     board.Color
}

// GetGroup finds all connected stones of the same color starting from (x, y).
func GetGroup(b *board.Board, x, y int) *Group {
	color := b.At(x, y)
	if color == board.Empty {
		return nil
	}

	visited := make(map[board.Point]bool)
	stones := []board.Point{}
	libertiesMap := make(map[board.Point]bool)

	var dfs func(int, int)
	dfs = func(cx, cy int) {
		p := board.Point{X: cx, Y: cy}
		if visited[p] {
			return
		}
		visited[p] = true
		stones = append(stones, p)

		adj := []board.Point{
			{X: cx + 1, Y: cy},
			{X: cx - 1, Y: cy},
			{X: cx, Y: cy + 1},
			{X: cx, Y: cy - 1},
		}

		for _, a := range adj {
			if b.IsOnBoard(a.X, a.Y) {
				if b.At(a.X, a.Y) == color {
					dfs(a.X, a.Y)
				} else if b.At(a.X, a.Y) == board.Empty {
					libertiesMap[a] = true
				}
			}
		}
	}

	dfs(x, y)

	liberties := make([]board.Point, 0, len(libertiesMap))
	for p := range libertiesMap {
		liberties = append(liberties, p)
	}

	return &Group{
		Stones:    stones,
		Liberties: liberties,
		Color:     color,
	}
}

// CountLiberties returns the number of liberties for the group at (x, y).
func CountLiberties(b *board.Board, x, y int) int {
	g := GetGroup(b, x, y)
	if g == nil {
		return 0
	}
	return len(g.Liberties)
}

// FindCapturedStones finds all opponent stones that would be captured by a move at (x, y) by color.
// This function assumes the move has already been tentatively placed (or simulated).
func FindCapturedStones(b *board.Board, x, y int, color board.Color) []board.Point {
	opponentColor := color.Opposite()
	captured := []board.Point{}
	visited := make(map[board.Point]bool)

	adj := []board.Point{
		{X: x + 1, Y: y},
		{X: x - 1, Y: y},
		{X: x, Y: y + 1},
		{X: x, Y: y - 1},
	}

	for _, a := range adj {
		if b.IsOnBoard(a.X, a.Y) && b.At(a.X, a.Y) == opponentColor {
			if !visited[a] {
				g := GetGroup(b, a.X, a.Y)
				if len(g.Liberties) == 0 {
					captured = append(captured, g.Stones...)
				}
				for _, s := range g.Stones {
					visited[s] = true
				}
			}
		}
	}

	return captured
}

// IsMoveValid checks if placing a stone of color at (x, y) is legal.
// It checks for:
// 1. Point is on board and empty.
// 2. Not suicide (unless it captures opponent stones).
// 3. TODO: Ko rule (needs game history/previous board state).
func IsMoveValid(b *board.Board, x, y int, color board.Color) bool {
	if !b.IsOnBoard(x, y) || b.At(x, y) != board.Empty {
		return false
	}

	// Tentatively place stone
	tempBoard := b.Copy()
	tempBoard.Set(x, y, color)

	// Check if this move captures anything
	captured := FindCapturedStones(tempBoard, x, y, color)
	if len(captured) > 0 {
		return true // Valid if it captures, even if it looks like suicide
	}

	// Check if the placed stone has liberties
	if CountLiberties(tempBoard, x, y) == 0 {
		return false // Suicide
	}

	return true
}
