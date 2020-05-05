package exception

import (
	"fmt"
	"os"
)

const (
	NormalError = iota
)

func CoralError(errType string, msg string) {
	fmt.Printf("* %s Error: %s", errType, msg)
	os.Exit(NormalError)
}
