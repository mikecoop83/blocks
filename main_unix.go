//go:build !js

package main

import (
	"log/slog"
	"os"
	"strconv"
)

func getGameIDFromParams() (uint64, error) {
	if len(os.Args) < 2 {
		return 0, nil
	}
	return strconv.ParseUint(os.Args[1], 10, 64)
}

func updateGameID(gameID uint64) {
	slog.Info("updating game id", "id", gameID)
}
