//go:build !js

package game

import "fmt"

func log(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
