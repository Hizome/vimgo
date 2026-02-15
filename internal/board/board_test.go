package board

import "testing"

func TestBoard(t *testing.T) {
	b := New(19)
	if b.Size != 19 {
		t.Errorf("expected size 19, got %d", b.Size)
	}
	b.Set(0, 0, Black)
	if b.At(0, 0) != Black {
		t.Errorf("expected Black at (0,0), got %v", b.At(0, 0))
	}
}
