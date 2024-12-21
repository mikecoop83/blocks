package main

import (
	"blocks/lib"
	"fmt"
)

func main() {
	b := lib.NewBoard()
	b.Place(lib.AllPieces[7], 0, 0)
	b.Place(lib.AllPieces[6], 5, 5)
	fmt.Println(b.String())
}
