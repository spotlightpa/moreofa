package errx_test

import (
	"fmt"

	"github.com/spotlightpa/moreofa/internal/errx"
)

func ExampleRecover() {
	maybePanic := func(throws bool) (err error) {
		defer errx.Recover(&err)

		if throws {
			panic("ahhh!")
		}
		return nil
	}

	err := maybePanic(false)
	fmt.Printf("error 1: %v\n", err)

	err = maybePanic(true)
	fmt.Printf("error 2: %v\n", err)

	// Output:
	// error 1: <nil>
	// error 2: panic: ahhh!
}
