package lib

import (
	_ "embed"
	"strings"
)

//go:embed pieces.txt
var piecesStr string

var AllPieces = parsePieces(piecesStr)

func parsePieces(input string) []Piece {
	pieces := make([]Piece, 0, 2)
	blocksStr := strings.Split(input, ".")

	for _, blockStr := range blocksStr {
		var shape [][]byte
		rows := strings.Split(blockStr, "\n")
		for _, row := range rows {
			var shapeRow []byte
			for _, char := range row {
				if char == '#' {
					shapeRow = append(shapeRow, 1)
				} else {
					shapeRow = append(shapeRow, 0)
				}
			}
			if len(shapeRow) > 0 {
				shape = append(shape, shapeRow)
			}
		}
		if len(shape) > 0 {
			pieces = append(pieces, Piece{Shape: shape})
		}
	}
	return pieces
}
