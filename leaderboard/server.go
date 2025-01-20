package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type LeaderboardEntry struct {
	GameID    uint64 `json:"game_id"`
	HighScore int    `json:"high_score"`
	Token     string `json:"token"`
}

var leaderboard = make(map[uint64]int)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var entry LeaderboardEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the token (this is a placeholder, implement actual validation)
	if entry.Token == "" {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Add the score to the leaderboard
	if currentHighScore, exists := leaderboard[entry.GameID]; !exists || entry.HighScore > currentHighScore {
		leaderboard[entry.GameID] = entry.HighScore
		slog.Info("New high score added", "game_id", entry.GameID, "high_score", entry.HighScore)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/add_score", Handler)
	slog.Info("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
