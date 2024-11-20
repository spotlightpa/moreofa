package errx_test

import (
	"errors"
	"fmt"

	"github.com/spotlightpa/moreofa/internal/errx"
)

func traceErr1(ok bool) (err error) {
	defer errx.Trace(&err)
	if !ok {
		return errors.New("oh no!")
	}
	return nil
}

func traceErr2(x, y int) (err error) {
	defer errx.Trace(&err)
	if x+y > 1 {
		return errors.New("uh oh!")
	}
	return nil
}

func ExampleTrace() {
	fmt.Println(traceErr1(true))
	fmt.Println(traceErr1(false))
	fmt.Println(traceErr2(1, -1))
	fmt.Println(traceErr2(1, 1))
	// Output:
	// <nil>
	// @github.com/spotlightpa/moreofa/internal/errx_test.traceErr1 (trace_example_test.go:13)
	// oh no!
	// <nil>
	// @github.com/spotlightpa/moreofa/internal/errx_test.traceErr2 (trace_example_test.go:21)
	// uh oh!
}
