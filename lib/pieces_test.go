package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePieces(t *testing.T) {
	require.Equal(t, Piece{Shape: [][]byte{{0, 1, 0}}}, parsePiece(" # "))
}
