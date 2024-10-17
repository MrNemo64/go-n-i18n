package util

import (
	"fmt"
	"os"
)

func Exit(code int, message string) {
	fmt.Print(message + "\n")
	os.Exit(code)
}
