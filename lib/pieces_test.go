package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePieces(t *testing.T) {
	piecesStr := `
###
 # 
.
###
###
###
`
	pieces := parsePieces(piecesStr)
	require.Equal(
		t,
		[]Piece{
			{
				Shape: [][]byte{
					{1, 1, 1},
					{0, 1, 0},
				},
			},
			{
				Shape: [][]byte{
					{1, 1, 1},
					{1, 1, 1},
					{1, 1, 1},
				},
			},
		}, pieces,
	)
}
