package lib

import (
	"strings"
)

type CellState int

const (
	Empty CellState = iota
	Occupied
)

type Board struct {
	Grid [8][8]CellState
}

type Piece struct {
	Shape [][]byte
}

func NewBoard() Board {
	return Board{}
}

func (b *Board) Place(p Piece, x, y int) bool {
	for i := range p.Shape {
		for j := range p.Shape[i] {
			if p.Shape[i][j] == 1 {
				if x+i >= 8 || y+j >= 8 || b.Grid[x+i][y+j] == 1 {
					return false
				}
			}
		}
	}
	for i := range p.Shape {
		for j := range p.Shape[i] {
			if p.Shape[i][j] == 1 {
				b.Grid[x+i][y+j] = 1
			}
		}
	}
	return true
}

func (b *Board) String() string {
	cellStateToIcon := map[CellState]string{
		Empty:    "⬜",
		Occupied: "🟦",
	}
	var sb strings.Builder
	for _, row := range b.Grid {
		for _, cell := range row {
			_, _ = sb.WriteString(cellStateToIcon[cell])
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
