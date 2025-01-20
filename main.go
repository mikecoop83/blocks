package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"syscall/js"

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

	js.Global().Set("wasmCallback", js.FuncOf(handleToken))

	// Run the game.
	err = ebiten.RunGame(game.New(gameID, updateGameID))
	if err != nil {
		panic(err)
	}
}
func handleToken(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 {
		token := args[0].String()
		fmt.Println("Received Google ID Token:", token)
		// Use the token to authenticate the player in your backend or game logic
	}
	return nil
}
