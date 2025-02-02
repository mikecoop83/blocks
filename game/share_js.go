//go:build js && wasm

package game

import (
	"syscall/js"
)

func copyToClipboard(text string) {
	navigator := js.Global().Get("navigator")
	if !navigator.Get("clipboard").IsUndefined() {
		navigator.Get("clipboard").Call("writeText", text)
	}
}

func getGameURL() string {
	location := js.Global().Get("window").Get("location")
	return location.Get("origin").String() + location.Get("pathname").String() + "?game=" + location.Get("search").String()[5:]
}
