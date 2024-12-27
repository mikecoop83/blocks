//go:build !js

package main

import "fmt"

func log(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
