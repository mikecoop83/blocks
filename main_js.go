//go:build js

package main

import (
	"errors"
	"log/slog"
	"net/url"
	"strconv"
	"syscall/js"
)

func getGameIDFromParams() (uint64, error) {
	query := js.Global().Get("window").Get("location").Get("search").String()
	if len(query) == 0 {
		return 0, errors.New("no query params")
	}
	values, err := url.ParseQuery(query[1:])
	if err != nil {
		return 0, err
	}
	slog.Info("query", "values", values)
	gameID := values.Get("game")
	if gameID == "" {
		return 0, errors.New("game param is empty")
	}
	return strconv.ParseUint(gameID, 16, 64)
}

func updateGameID(gameID uint64) {
	query := url.Values{}
	query.Set("game", strconv.FormatUint(gameID, 16))
	js.Global().Get("window").Get("history").Call("replaceState", nil, "", "?"+query.Encode())
}
