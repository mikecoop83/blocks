package main

import (
	"fmt"
	"math/rand"

	"github.com/mikecoop83/blocks/lib"

	"github.com/eiannone/keyboard"
)

func main() {
	err := mainWithError()
	if err != nil {
		fmt.Println(err)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func mainWithError() error {

	if err := keyboard.Open(); err != nil {
		return err
	}
	defer func() {
		_ = keyboard.Close()
	}()

	clearScreen()
	b := lib.NewBoard()

outer:
	for {
		piece := lib.AllPieces[rand.Intn(len(lib.AllPieces))]
		rotateTimes := rand.Intn(4)
		for i := 0; i < rotateTimes; i++ {
			piece = piece.Rotate()
		}
		loc := lib.Location{}
	inner:
		for {
			newGrid, _ := b.AddPiece(
				lib.PieceLocation{
					Piece: piece,
					Loc:   loc,
				},
				true,
			)
			clearScreen()
			fmt.Println(newGrid)

			char, key, err := keyboard.GetKey()
			if err != nil {
				return err
			}

			if char == 'q' {
				fmt.Println("Quitting...")
				break outer
			} else if char == 'z' {
				b.Undo()
				break
			} else if char == 'n' {
				break
			}

			newLoc := loc
			switch key {
			case keyboard.KeyArrowUp:
				if newLoc.X > 0 {
					newLoc.X--
				}
			case keyboard.KeyArrowDown:
				if newLoc.X < 7 {
					newLoc.X++
				}
			case keyboard.KeyArrowLeft:
				if newLoc.Y > 0 {
					newLoc.Y--
				}
			case keyboard.KeyArrowRight:
				if newLoc.Y < 7 {
					newLoc.Y++
				}
			case keyboard.KeyEnter:
				_, newValid := b.AddPiece(
					lib.PieceLocation{
						Piece: piece,
						Loc:   newLoc,
					},
					false,
				)
				if !newValid {
					continue
				}
				break inner
			default:
				// ignore
			}

			if b.ValidatePiece(
				lib.PieceLocation{
					Piece: piece,
					Loc:   newLoc,
				}, true,
			) {
				loc = newLoc
			}
		}
	}
	return nil
}
