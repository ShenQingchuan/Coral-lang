package exception

import (
	"errors"
	"fmt"
	"os"
)

type CoralError struct {
	Err     error
	ErrEnum int
}

func NewCoralError(prefixDescription string, msg string, errEnum int) *CoralError {
	return &CoralError{
		errors.New("* " + prefixDescription + " Error: " + msg),
		errEnum,
	}
}

func CoralErrorCrashHandler(c *CoralError) {
	fmt.Println(c.Err)
	fmt.Printf("* Error code: %d", c.ErrEnum)
	os.Exit(c.ErrEnum)
}
