package game

import (
	"log/slog"
	"strconv"

	"github.com/mikecoop83/blocks/persist"
)

func maybeGetHighScore() int64 {
	highScoreStr, err := persist.Load("highscore")
	if err != nil {
		slog.Error("failed to load high score: %v", err)
		return 0
	}
	if highScoreStr == "" {
		return 0
	}
	highScore, err := strconv.ParseInt(highScoreStr, 10, 64)
	if err != nil {
		slog.Error("failed to parse high score: %v", err)
		return 0
	}
	return highScore
}

func maybeUpdateHighScore(highScore int64) {
	err := persist.Store("highscore", strconv.FormatInt(highScore, 10))
	if err != nil {
		slog.Error("failed to save high score: %v", err)
	}
}
