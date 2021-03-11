package exception

import (
	. "coral-lang/src/utils"
	"errors"
	"fmt"
	"os"
	"strings"
)

type CoralCompileError struct {
	Err     error
	ErrEnum int
}

func NewCoralError(prefixDescription string, msg string, errEnum int) *CoralCompileError {
	return &CoralCompileError{
		errors.New("\n* " + Red(prefixDescription+" Error: ") + msg),
		errEnum,
	}
}

func CoralErrorCrashHandler(c *CoralCompileError) {
	fmt.Println(c.Err)
	fmt.Println(Cyan(fmt.Sprintf("* Error code: %d", c.ErrEnum)))
	os.Exit(c.ErrEnum)
}

func CoralCompileWarning(msg string) {
	fmt.Println("\n" + Yellow("* Warning: "))
	for _, str := range strings.Split(msg, "\n") {
		fmt.Println("\t" + White(str))
	}
}
