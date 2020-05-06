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

func NewCoralError(errKind string, msg string, errEnum int) *CoralError {
	return &CoralError{
		errors.New("* " + errKind + " Error: " + msg),
		errEnum,
	}
}

func CoralErrorHandler(c *CoralError) {
	fmt.Println(c.Err)
	os.Exit(c.ErrEnum)
}
