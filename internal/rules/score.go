package rules

import (
	"strings"

	"github.com/vimgo/vimgo/internal/board"
)

type Score struct {
	Black float64
	White float64
}

// CountScore calculates the score for the current board state.
// It assumes all stones on the board are alive.
// method: "chinese" (area) or "japanese" (territory)
func CountScore(b *board.Board, method string, blackCaptures, whiteCaptures int, komi float64) Score {
	normalized := strings.ToLower(strings.TrimSpace(method))
	if normalized != "japanese" {
		normalized = "chinese"
	}

	blackScore := 0.0
	whiteScore := komi

	// Territory counting (Japanese) adds captures to score.
	// Area counting (Chinese) counts stones + territory.
	if normalized == "japanese" {
		blackScore += float64(blackCaptures)
		whiteScore += float64(whiteCaptures)
	}

	visited := make([]bool, b.Size*b.Size)

	for y := 0; y < b.Size; y++ {
		for x := 0; x < b.Size; x++ {
			c := b.At(x, y)

			// Chinese: count stones
			if normalized == "chinese" {
				if c == board.Black {
					blackScore++
				} else if c == board.White {
					whiteScore++
				}
			}

			// Territory detection (empty points)
			idx := y*b.Size + x
			if c == board.Empty && !visited[idx] {
				points, owner := getTerritory(b, x, y, visited)
				if owner == board.Black {
					blackScore += float64(len(points))
				} else if owner == board.White {
					whiteScore += float64(len(points))
				}
			}
		}
	}

	return Score{Black: blackScore, White: whiteScore}
}

// getTerritory performs flood fill to find connected empty points and determines ownership.
// Returns the list of points and the owner (Black, White, or Empty if shared/dame).
func getTerritory(b *board.Board, startX, startY int, visited []bool) ([]board.Point, board.Color) {
	points := []board.Point{}
	touchedBlack := false
	touchedWhite := false

	queue := []board.Point{{X: startX, Y: startY}}
	visited[startY*b.Size+startX] = true

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		points = append(points, p)

		adj := []board.Point{
			{X: p.X + 1, Y: p.Y},
			{X: p.X - 1, Y: p.Y},
			{X: p.X, Y: p.Y + 1},
			{X: p.X, Y: p.Y - 1},
		}

		for _, a := range adj {
			if !b.IsOnBoard(a.X, a.Y) {
				continue
			}

			pixel := b.At(a.X, a.Y)
			if pixel == board.Empty {
				idx := a.Y*b.Size + a.X
				if !visited[idx] {
					visited[idx] = true
					queue = append(queue, a)
				}
			} else if pixel == board.Black {
				touchedBlack = true
			} else if pixel == board.White {
				touchedWhite = true
			}
		}
	}

	if touchedBlack && !touchedWhite {
		return points, board.Black
	}
	if touchedWhite && !touchedBlack {
		return points, board.White
	}
	return points, board.Empty // Dame or neutral
}
