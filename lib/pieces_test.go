package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePieces(t *testing.T) {
	require.Equal(t, Piece{Shape: [][]bool{{false, true, false}, {true, true, true}}}, parsePiece(" # \n###"))
}
