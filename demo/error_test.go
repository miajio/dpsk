package demo_test

import (
	"fmt"
	"testing"

	"github.com/miajio/dpsk/errors"
)

func TestError(t *testing.T) {
	err := errors.NewF("test error")
	fmt.Println(err)

	err = errors.NewCodeErrorF(500, "test error")
	fmt.Println(err)

	errors.ErrPrintTypeDefault = errors.ErrPrintTypeJson
	fmt.Println(err)

	codeErr := errors.ReadCodeError(err)
	fmt.Println(codeErr.Code)
}
