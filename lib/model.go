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
	Unchosen
	Hovering
)

const BoardSize = 8

type Grid [BoardSize][BoardSize]CellState

func (g Grid) String() string {
	cellStateToIcon := map[CellState]string{
		Empty:    "e",
		Pending:  "p",
		Invalid:  "i",
		FullLine: "f",
		Occupied: "o",
		Unchosen: "u",
		Hovering: "h",
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

type Board struct {
	gridHistory *Stack[Grid]
}

type Piece struct {
	Shape [][]bool
}

func (p Piece) Height() int {
	return len(p.Shape)
}

func (p Piece) Width() int {
	return len(p.Shape[0])
}

func (p Piece) Rotate() Piece {
	rotated := Piece{
		Shape: make([][]bool, len(p.Shape[0])),
	}
	for c := range p.Shape[0] {
		rotated.Shape[c] = make([]bool, len(p.Shape))
	}
	for r := range p.Shape {
		for c := range p.Shape[r] {
			rotated.Shape[c][len(p.Shape)-1-r] = p.Shape[r][c]
		}
	}
	return rotated
}

func (p Piece) NumBlocks() int {
	var numBlocks int
	for r := range p.Shape {
		for c := range p.Shape[r] {
			if p.Shape[r][c] {
				numBlocks++
			}
		}
	}
	return numBlocks
}

func NewBoard() Board {
	gridHistory := NewStack[Grid]()
	gridHistory.Push(Grid{})
	return Board{
		gridHistory: gridHistory,
	}
}

type Location struct {
	C, R int
}

type PieceLocation struct {
	Piece Piece
	Loc   Location
}

func (b *Board) Clear() {
	b.gridHistory = NewStack[Grid]()
}

func (b *Board) CanPlacePiece(piece Piece) bool {
	for r := range piece.Shape {
		for c := range piece.Shape[r] {
			loc := Location{C: c, R: r}
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
	if loc.C < 0 || loc.R < 0 {
		return false
	}
	if loc.C+piece.Width() > BoardSize || loc.R+piece.Height() > BoardSize {
		return false
	}
	if allowPieceOverlap {
		return true
	}
	for r := range piece.Shape {
		for c := range piece.Shape[r] {
			if piece.Shape[r][c] {
				if grid[loc.R+r][loc.C+c] == Occupied {
					return false
				}
			}
		}
	}
	return true
}

func (b *Board) AddPiece(pieceLoc PieceLocation, pending bool) (Grid, int, bool) {
	if !b.ValidatePiece(pieceLoc, pending) {
		return b.GetGrid(), 0, false
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
	for r := range piece.Shape {
		for c := range piece.Shape[r] {
			if !piece.Shape[r][c] {
				continue
			}
			if grid[loc.R+r][loc.C+c] == Occupied {
				grid[loc.R+r][loc.C+c] = Invalid
				anyInvalid = true
				continue
			}
			grid[loc.R+r][loc.C+c] = newPieceState
		}
	}
	if anyInvalid {
		return grid, 0, false
	}

	var clearedLines int
	// Find full horizontal lines
	for r := range BoardSize {
		full := true
		for c := range BoardSize {
			if grid[r][c] != Occupied && grid[r][c] != Pending {
				full = false
				break
			}
		}
		if !full {
			continue
		}
		clearedLines++
		for c := range BoardSize {
			if grid[r][c] != Pending {
				grid[r][c] = newFullLineState
			}
		}
	}
	// Find full vertical lines
	for c := range BoardSize {
		full := true
		for r := range BoardSize {
			if grid[r][c] != Occupied && grid[r][c] != Pending {
				full = false
				break
			}
		}
		if !full {
			continue
		}
		clearedLines++
		for r := range BoardSize {
			if grid[r][c] != Pending {
				grid[r][c] = newFullLineState
			}
		}
	}
	if !pending {
		b.gridHistory.Push(grid)
	}
	return grid, clearedLines, true
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
