package rules

import (
	"testing"

	"github.com/vimgo/vimgo/internal/board"
)

func TestCountScoreChinese(t *testing.T) {
	b := board.New(5)
	// Black corner shape, 1-point corner territory at (0,0)
	b.Set(0, 1, board.Black)
	b.Set(1, 0, board.Black)
	// White corner shape, 1-point corner territory at (4,4)
	b.Set(4, 3, board.White)
	b.Set(3, 4, board.White)

	score := CountScore(b, "chinese", 0, 0, 7.5)

	if score.Black != 3.0 {
		t.Fatalf("black score mismatch: got %.1f want 3.0", score.Black)
	}
	if score.White != 10.5 {
		t.Fatalf("white score mismatch: got %.1f want 10.5", score.White)
	}
}

func TestCountScoreJapanese(t *testing.T) {
	b := board.New(5)
	b.Set(0, 1, board.Black)
	b.Set(1, 0, board.Black)
	b.Set(4, 3, board.White)
	b.Set(3, 4, board.White)

	score := CountScore(b, "japanese", 2, 1, 6.5)

	if score.Black != 3.0 {
		t.Fatalf("black score mismatch: got %.1f want 3.0", score.Black)
	}
	if score.White != 8.5 {
		t.Fatalf("white score mismatch: got %.1f want 8.5", score.White)
	}
}

func TestCountScoreNeutralTerritoryNotCounted(t *testing.T) {
	b := board.New(3)
	b.Set(0, 1, board.Black)
	b.Set(2, 1, board.White)

	score := CountScore(b, "japanese", 0, 0, 0)

	if score.Black != 0.0 {
		t.Fatalf("black score mismatch: got %.1f want 0.0", score.Black)
	}
	if score.White != 0.0 {
		t.Fatalf("white score mismatch: got %.1f want 0.0", score.White)
	}
}

func TestCountScoreUnknownMethodFallsBackToChinese(t *testing.T) {
	b := board.New(3)
	b.Set(0, 1, board.Black)
	b.Set(2, 1, board.White)

	score := CountScore(b, "bad-input", 99, 99, 0)

	if score.Black != 1.0 {
		t.Fatalf("black score mismatch: got %.1f want 1.0", score.Black)
	}
	if score.White != 1.0 {
		t.Fatalf("white score mismatch: got %.1f want 1.0", score.White)
	}
}
