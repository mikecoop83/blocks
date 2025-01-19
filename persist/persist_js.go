//go:build js

package persist

import "syscall/js"

func Store(key string, value string) error {
	localStorage := js.Global().Get("localStorage")
	localStorage.Call("setItem", key, value)
	return nil
}

func Load(key string) (string, error) {
	localStorage := js.Global().Get("localStorage")
	val := localStorage.Call("getItem", key)
	if val.IsNull() {
		return "", nil
	}
	return val.String(), nil
}
