//go:build !(js && wasm)

package game

func getGameURL() string {
	return ""
}

func copyToClipboard(text string) {
	// No-op for non-JS platforms
}
