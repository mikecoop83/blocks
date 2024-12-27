package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikecoop83/blocks/game"
)

func main() {
	// Set the window title.
	ebiten.SetWindowTitle("Blocks")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(game.WindowWidth, game.WindowHeight)

	// Run the game.
	err := ebiten.RunGame(game.New())
	if err != nil {
		panic(err)
	}
}
