package lib

import (
	"strings"
)

type CellState int

const (
	Empty CellState = iota
	Pending
	Invalid
	FullLine
	Occupied
)

const BoardSize = 8

type Grid [BoardSize][BoardSize]CellState

type Board struct {
	gridHistory *Stack[Grid]
}

type Piece struct {
	Shape [][]bool
}

func (p Piece) Rotate() Piece {
	rotated := Piece{
		Shape: make([][]bool, len(p.Shape[0])),
	}
	for i := range p.Shape[0] {
		rotated.Shape[i] = make([]bool, len(p.Shape))
	}
	for i := range p.Shape {
		for j := range p.Shape[i] {
			rotated.Shape[j][len(p.Shape)-1-i] = p.Shape[i][j]
		}
	}
	return rotated
}

func NewBoard() Board {
	gridHistory := NewStack[Grid]()
	gridHistory.Push(Grid{})
	return Board{
		gridHistory: gridHistory,
	}
}

type Location struct {
	X, Y int
}

type PieceLocation struct {
	Piece Piece
	Loc   Location
}

func (b *Board) Clear() {
	b.gridHistory = NewStack[Grid]()
}

func (b *Board) CanPlacePiece(piece Piece) bool {
	for i := range piece.Shape {
		for j := range piece.Shape[i] {
			loc := Location{X: i, Y: j}
			if b.ValidatePiece(PieceLocation{Piece: piece, Loc: loc}, false) {
				return true
			}
		}
	}
	return false
}

func (b *Board) ValidatePiece(
	pieceLoc PieceLocation,
	allowPieceOverlap bool,
) bool {
	grid, ok := b.gridHistory.Peek()
	if !ok {
		panic("no grid history")
	}
	loc := pieceLoc.Loc
	piece := pieceLoc.Piece
	for i := range piece.Shape {
		for j := range piece.Shape[i] {
			if piece.Shape[i][j] {
				if loc.X+i >= BoardSize || loc.Y+j >= BoardSize {
					return false
				}
				if allowPieceOverlap {
					continue
				}
				if grid[loc.X+i][loc.Y+j] == Occupied {
					return false
				}
			}
		}
	}
	return true
}

func (b *Board) AddPiece(pieceLoc PieceLocation, pending bool) (Grid, bool) {
	if !b.ValidatePiece(pieceLoc, pending) {
		return b.GetGrid(), false
	}
	grid, ok := b.gridHistory.Peek()
	if !ok {
		panic("no grid history")
	}
	piece := pieceLoc.Piece
	loc := pieceLoc.Loc
	anyInvalid := false

	newPieceState := Occupied
	newFullLineState := Empty
	if pending {
		newPieceState = Pending
		newFullLineState = FullLine
	}

	// Visualize the piece on the board
	for i := range piece.Shape {
		for j := range piece.Shape[i] {
			if !piece.Shape[i][j] {
				continue
			}
			if grid[loc.X+i][loc.Y+j] == Occupied {
				grid[loc.X+i][loc.Y+j] = Invalid
				anyInvalid = true
				continue
			}
			grid[loc.X+i][loc.Y+j] = newPieceState
		}
	}
	if anyInvalid {
		return grid, false
	}

	// Find full horizontal lines
	for i := 0; i < BoardSize; i++ {
		full := true
		for j := 0; j < BoardSize; j++ {
			if grid[i][j] != Occupied && grid[i][j] != Pending {
				full = false
				break
			}
		}
		if !full {
			continue
		}
		for j := 0; j < BoardSize; j++ {
			if grid[i][j] != Pending {
				grid[i][j] = newFullLineState
			}
		}
	}
	// Find full vertical lines
	for j := 0; j < BoardSize; j++ {
		full := true
		for i := 0; i < BoardSize; i++ {
			if grid[i][j] != Occupied && grid[i][j] != Pending {
				full = false
				break
			}
		}
		if !full {
			continue
		}
		for i := 0; i < BoardSize; i++ {
			if grid[i][j] != Pending {
				grid[i][j] = newFullLineState
			}
		}
	}
	if !pending {
		b.gridHistory.Push(grid)
	}
	return grid, true
}

func (b *Board) Undo() bool {
	if b.gridHistory.Len() == 1 {
		return false
	}
	b.gridHistory.Pop()
	return true
}

func (b *Board) GetGrid() Grid {
	grid, ok := b.gridHistory.Peek()
	if !ok {
		panic("no grid history")
	}
	return grid
}

func (g Grid) String() string {
	cellStateToIcon := map[CellState]string{
		Empty:    "⬜",
		Pending:  "🟩",
		Invalid:  "🟥",
		FullLine: "🟧",
		Occupied: "🟦",
	}
	var sb strings.Builder
	for _, row := range g {
		for _, cell := range row {
			_, _ = sb.WriteString(cellStateToIcon[cell])
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
