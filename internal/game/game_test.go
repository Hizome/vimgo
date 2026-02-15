package game

import (
	"github.com/vimgo/vimgo/internal/board"
	"testing"
)

func TestGame_Move(t *testing.T) {
	g := NewGame(9)

	// Test basic moves
	if err := g.Move(4, 4); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if g.Board.At(4, 4) != board.Black {
		t.Errorf("expected Black at (4,4), got %v", g.Board.At(4, 4))
	}
	if g.CurrentPlayer != board.White {
		t.Errorf("expected player to be White, got %v", g.CurrentPlayer)
	}

	// Test capture
	g = NewGame(9)
	g.Move(0, 1) // Black
	g.Move(0, 0) // White
	g.Move(1, 0) // Black
	// White at (0,0) is captured if Black plays at (0,1) and (1,0)
	// Wait, (0,0) liberties are (0,1) and (1,0).
	// If Black occupies both, White (0,0) is captured.
	if g.Board.At(0, 0) != board.Empty {
		t.Errorf("expected (0,0) to be empty after capture, got %v", g.Board.At(0, 0))
	}
	if g.BlackCaptures != 1 {
		t.Errorf("expected Black to have 1 capture, got %d", g.BlackCaptures)
	}
}

func TestGame_Ko(t *testing.T) {
	g := NewGame(9)
	// Setup Ko situation
	// B: (1,0), (0,1), (2,1), (1,2)
	// W: (0,0), (1,1) -> wait, simple Ko:
	// . B .
	// B W B
	// . B .
	// White plays at (1,1) then Black plays at (1,1) capturing White

	g.Move(1, 0) // B
	g.Move(0, 0) // W
	g.Move(0, 1) // B
	g.Move(2, 0) // W
	g.Move(1, 1) // B
	g.Move(1, 0) // W - this should capture B at (1,1)

	// Now B at (1,1) is empty. B tries to capture W back at (1,1) immediately.
	err := g.Move(1, 1)
	if err == nil {
		t.Error("expected Ko violation error, got nil")
	}
}

func TestGame_UndoRestoresCapturesAndMetadata(t *testing.T) {
	g := NewGame(9)
	if err := g.Move(0, 1); err != nil { // B
		t.Fatalf("unexpected move error: %v", err)
	}
	if err := g.Move(0, 0); err != nil { // W
		t.Fatalf("unexpected move error: %v", err)
	}
	if err := g.Move(1, 0); err != nil { // B captures W(0,0)
		t.Fatalf("unexpected move error: %v", err)
	}

	if g.BlackCaptures != 1 {
		t.Fatalf("expected BlackCaptures=1 before undo, got %d", g.BlackCaptures)
	}
	if len(g.Moves) != 3 {
		t.Fatalf("expected 3 moves before undo, got %d", len(g.Moves))
	}

	if err := g.Undo(); err != nil {
		t.Fatalf("unexpected undo error: %v", err)
	}

	if g.Board.At(0, 0) != board.White {
		t.Fatalf("expected white stone restored at (0,0), got %v", g.Board.At(0, 0))
	}
	if g.Board.At(1, 0) != board.Empty {
		t.Fatalf("expected (1,0) to be empty after undo, got %v", g.Board.At(1, 0))
	}
	if g.BlackCaptures != 0 {
		t.Fatalf("expected BlackCaptures restored to 0, got %d", g.BlackCaptures)
	}
	if g.WhiteCaptures != 0 {
		t.Fatalf("expected WhiteCaptures restored to 0, got %d", g.WhiteCaptures)
	}
	if g.CurrentPlayer != board.Black {
		t.Fatalf("expected current player restored to Black, got %v", g.CurrentPlayer)
	}
	if len(g.Moves) != 2 {
		t.Fatalf("expected moves length restored to 2, got %d", len(g.Moves))
	}
	if g.LastMove == nil || g.LastMove.X != 0 || g.LastMove.Y != 0 {
		t.Fatalf("expected last move restored to (0,0), got %+v", g.LastMove)
	}
}
