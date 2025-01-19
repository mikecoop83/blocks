package main

import (
	"log/slog"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mikecoop83/blocks/game"
)

func main() {
	// Set the window title.
	ebiten.SetWindowTitle("Blocks")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(game.WindowWidth, game.WindowHeight)

	gameID, err := getGameIDFromParams()
	if err != nil {
		slog.Error("unable to parse game ID", "error", err)
	}
	if gameID == 0 {
		slog.Info("no game ID found, generating a new one")
		gameID = rand.Uint64()
	}

	// Run the game.
	err = ebiten.RunGame(game.New(gameID, updateGameID))
	if err != nil {
		panic(err)
	}
}
