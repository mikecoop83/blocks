package lib

import (
	_ "embed"
	"math/rand"
	"strings"
)

var AllPieces []Piece

func init() {
	AllPieces = []Piece{
		/* 0 */ parsePiece("#"),
		/* 1 */ parsePiece("##"),
		/* 2 */ parsePiece("###"),
		/* 3 */ parsePiece("####"),
		/* 4 */ parsePiece("#####"),
		/* 5 */ parsePiece("##\n##"),
		/* 6 */ parsePiece("###\n###\n###"),
		/* 7 */ parsePiece("###\n # "),
		/* 8 */ parsePiece("# \n# \n##"),
		/* 9 */ parsePiece(" #\n #\n##"),
		/* 10 */ parsePiece(" ##\n## "),
		/* 11 */ parsePiece("## \n ##"),
		/* 12 */ parsePiece("##\n #"),
		/* 13 */ parsePiece("# \n##"),
		/* 14 */ parsePiece("#  \n # \n  #"),
		/* 15 */ parsePiece("# \n #\n"),
		/* 16 */ parsePiece("###\n###"),
		/* 17 */ parsePiece("#  \n#  \n###"),
		/* 18 */ parsePiece("  #\n  #\n###"),
	}
}

func parsePiece(pieceStr string) Piece {
	rows := strings.Split(pieceStr, "\n")
	shape := make([][]bool, 0, len(rows))
	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		pieceRow := make([]bool, len(row))
		for i, c := range row {
			if c != ' ' {
				pieceRow[i] = true
			}
		}
		shape = append(shape, pieceRow)
	}
	return Piece{
		Shape: shape,
	}
}

func RandomRotatedPiece() Piece {
	randPieceIdx := rand.Intn(len(AllPieces))
	piece := AllPieces[randPieceIdx]
	rotateTimes := rand.Intn(4)
	for i := 0; i < rotateTimes; i++ {
		piece = piece.Rotate()
	}
	return piece
}
