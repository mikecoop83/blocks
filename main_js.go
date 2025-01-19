//go:build js

package main

import (
	"net/url"
	"strconv"
	"syscall/js"
)

func getGameIDFromParams() (uint64, error) {
	query := js.Global().Get("window").Get("location").Get("search").String()
	values, err := url.ParseQuery(query)
	if err != nil {
		return 0, err
	}
	gameID := values.Get("game")
	if gameID == "" {
		return 0, nil
	}
	return strconv.ParseUint(gameID, 10, 64)
}

func updateGameID(gameID uint64) {
	query := url.Values{}
	query.Set("game", strconv.FormatUint(gameID, 10))
	js.Global().Get("window").Get("history").Call("replaceState", nil, "", "?"+query.Encode())
}
