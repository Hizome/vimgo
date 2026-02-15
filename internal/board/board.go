package board



type Color int

const (
	Empty Color = iota
	Black
	White
)

func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "Empty"
	}
}

func (c Color) Opposite() Color {
	if c == Black {
		return White
	}
	if c == White {
		return Black
	}
	return Empty
}

type Point struct {
	X, Y int
}

type Board struct {
	Size int
	Grid []Color
}

func New(size int) *Board {
	return &Board{
		Size: size,
		Grid: make([]Color, size*size),
	}
}

func (b *Board) At(x, y int) Color {
	if x < 0 || x >= b.Size || y < 0 || y >= b.Size {
		return Empty // Or panic/error? For now Empty is safe for some checks, but maybe not for captures.
	}
	return b.Grid[y*b.Size+x]
}

func (b *Board) Set(x, y int, c Color) {
	if x >= 0 && x < b.Size && y >= 0 && y < b.Size {
		b.Grid[y*b.Size+x] = c
	}
}

func (b *Board) IsOnBoard(x, y int) bool {
	return x >= 0 && x < b.Size && y >= 0 && y < b.Size
}

func (b *Board) Clear() {
	for i := range b.Grid {
		b.Grid[i] = Empty
	}
}

func (b *Board) Copy() *Board {
	newB := New(b.Size)
	copy(newB.Grid, b.Grid)
	return newB
}
