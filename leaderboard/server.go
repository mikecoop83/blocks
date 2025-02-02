package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/api/idtoken"
)

type Request struct {
	GameID    uint64 `json:"game_id"`
	HighScore int    `json:"high_score"`
	Token     string `json:"token"`
}

type Score struct {
	GameID uint64 `json:"game_id"`
	Score  int    `json:"score"`
	User   User   `json:"user"`
	Time   int64  `json:"time"`
}

type User struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func validateToken(ctx context.Context, token string) (User, error) {
	if token == "" {
		return User{}, fmt.Errorf("empty token")
	}
	clientID := os.Getenv("CLIENT_ID")
	payload, err := idtoken.Validate(ctx, token, clientID)
	if err != nil {
		return User{}, err
	}
	expiresTime := time.Unix(payload.Expires, 0)
	user := User{
		Email:   payload.Claims["email"].(string),
		Name:    payload.Claims["name"].(string),
		Picture: payload.Claims["picture"].(string),
	}
	slog.Info(
		"Token validated", "email", user.Email, "expires", expiresTime, "issued", payload.IssuedAt, "claims",
		payload.Claims,
	)
	err = createUser(ctx, user)
	return user, err
}

func Handler(w http.ResponseWriter, r *http.Request) {
	writeLeaderboard := func() {
		w.Header().Set("Content-Type", "application/json")
		scores, err := getTopScores(r.Context(), 10)
		if err != nil {
			slog.Error("Failed to get scores", "error", err)
			http.Error(w, "Failed to get scores", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(scores)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
	if r.Method == http.MethodGet {
		writeLeaderboard()
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var entry Request
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := validateToken(r.Context(), entry.Token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %s", err.Error()), http.StatusUnauthorized)
		return
	}

	// Add the score to the leaderboard
	score := Score{
		GameID: entry.GameID,
		Score:  entry.HighScore,
		User:   user,
		Time:   time.Now().UTC().Unix(),
	}
	err = addScore(r.Context(), score)
	if err != nil {
		http.Error(w, "Failed to add score", http.StatusInternalServerError)
		return
	}
	writeLeaderboard()
}

func createUser(ctx context.Context, user User) error {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()
	query := `INSERT INTO public.users (email, name, picture)
VALUES ($1, $2, $3)
ON CONFLICT (email) DO UPDATE SET
    name = EXCLUDED.name,
    picture = EXCLUDED.picture;`
	_, err = db.ExecContext(ctx, query, user.Email, user.Name, user.Picture)
	if err != nil {
		return err
	}
	slog.Info("Inserted user", "email", user.Email, "name", user.Name, "picture", user.Picture)
	return nil
}

func GetUserByEmail(ctx context.Context, email string) (User, error) {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return User{}, err
	}
	defer func() {
		_ = db.Close()
	}()
	query := `SELECT email, name, picture FROM users WHERE email = $1`
	row := db.QueryRowContext(ctx, query, email)
	var user User
	err = row.Scan(&user.Email, &user.Name, &user.Picture)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func addScore(ctx context.Context, score Score) error {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()
	query := `INSERT INTO public.scores (score, game_id, user_email) VALUES ($1, $2, $3)`
	_, err = db.ExecContext(ctx, query, score.Score, score.GameID, score.User.Email)
	return err
}

func getTopScores(ctx context.Context, limit int) ([]Score, error) {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = db.Close()
	}()
	query := `SELECT s.score, s.created, s.game_id, s.user_email, u.name, u.picture
FROM scores s
LEFT JOIN users u
  on s.user_email = u.email
ORDER BY score DESC, created ASC
LIMIT $1;`
	rows, err := db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	scores := make([]Score, 0, limit)
	for rows.Next() {
		var score Score
		var created time.Time
		err = rows.Scan(
			&score.Score,
			&created,
			&score.GameID,
			&score.User.Email,
			&score.User.Name,
			&score.User.Picture,
		)
		if err != nil {
			return nil, err
		}
		score.Time = created.UTC().Unix()
		scores = append(scores, score)
	}
	return scores, nil
}
