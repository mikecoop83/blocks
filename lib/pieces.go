package lib

import (
	_ "embed"
	"strings"
)

var AllPieces []Piece

func init() {
	AllPieces = []Piece{
		parsePiece("#"),
		parsePiece("##"),
		parsePiece("###"),
		parsePiece("####"),
		parsePiece("#####"),
		parsePiece("##\n##"),
		parsePiece("###\n###\n###"),
		parsePiece("###\n #"),
		parsePiece("#\n#\n##"),
		parsePiece(" #\n #\n##"),
		parsePiece(" ##\n##"),
		parsePiece("##\n ##"),
		parsePiece("##\n #"),
		parsePiece("#\n##"),
	}
}
func parsePiece(pieceStr string) Piece {
	rows := strings.Split(pieceStr, "\n")
	shape := make([][]byte, 0, len(rows))
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		pieceRow := make([]byte, 0, len(row))
		for _, c := range row {
			if c == ' ' {
				pieceRow = append(pieceRow, 0)
			} else {
				pieceRow = append(pieceRow, 1)
			}
		}
		shape = append(shape, pieceRow)
	}
	return Piece{
		Shape: shape,
	}
}
