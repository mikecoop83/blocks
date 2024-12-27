//go:build js

package game

import (
	"fmt"
	"syscall/js"
)

func log(msg string, args ...interface{}) {
	js.Global().Get("console").Call("log", fmt.Sprintf(msg, args...))
}
