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

const boardSize = 8

type Grid [boardSize][boardSize]CellState

type Board struct {
	gridHistory *Stack[Grid]
}

type Piece struct {
	Shape [][]bool
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
				if loc.X+i >= boardSize || loc.Y+j >= boardSize {
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
		return Grid{}, false
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
	for i := 0; i < boardSize; i++ {
		full := true
		for j := 0; j < boardSize; j++ {
			if grid[i][j] != Occupied && grid[i][j] != Pending {
				full = false
				break
			}
		}
		if !full {
			continue
		}
		for j := 0; j < boardSize; j++ {
			grid[i][j] = newFullLineState
		}
	}
	// Find full vertical lines
	for j := 0; j < boardSize; j++ {
		full := true
		for i := 0; i < boardSize; i++ {
			if grid[i][j] != Occupied && grid[i][j] != Pending {
				full = false
				break
			}
		}
		if !full {
			continue
		}
		for i := 0; i < boardSize; i++ {
			grid[i][j] = newFullLineState
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
