package errx_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/carlmjohnson/be"
	"github.com/spotlightpa/moreofa/internal/errx"
)

func ExampleValidator() {
	var v errx.Validator
	v.AddIf(2 > 1, "heads", "Two are better than one.")
	v.AddIf(true, "heads", "I win, tails you lose.")
	err := v.Err()

	fmt.Println(errx.StatusCode(err))
	for field, msgs := range errx.ValidationErrors(err) {
		for _, msg := range msgs {
			fmt.Println(field, "=", msg)
		}
	}
	// Output:
	// 400
	// heads = Two are better than one.
	// heads = I win, tails you lose.
}

func ExampleValidator_AddIfUnset() {
	var v errx.Validator
	x, err := strconv.Atoi("hello")
	v.AddIf(err != nil, "x", "Could not parse x.")
	v.AddIf(x < 1, "x", "X must be positive.")

	y, err := strconv.Atoi("hello")
	v.AddIf(err != nil, "y", "Could not parse y.")
	v.AddIfUnset(y < 1, "y", "Y must be positive.")
	fmt.Println(v.Err())
	// Output:
	// validation error: x=Could not parse x. x=X must be positive. y=Could not parse y.
}

func TestValidator(t *testing.T) {
	var v1 errx.Validator
	v1.AddIf(2 > 1, "heads", "Two are better than one.")
	v1.AddIf(true, "heads", "I win, tails you lose.")
	err := v1.Err()
	be.Nonzero(t, err)
	be.False(t, v1.Valid())
	fields := errx.ValidationErrors(err)
	be.Equal(t, 1, len(fields))
	be.Equal(t, 2, len(fields["heads"]))
	be.Equal(t, http.StatusBadRequest, errx.StatusCode(err))

	var v2 errx.Validator
	v2.AddIf(2 < 1, "heads", "One is the loneliest number.")
	v2.AddIf(false, "heads", "I win, tails you lose.")
	err = v2.Err()
	be.True(t, v2.Valid())
	be.NilErr(t, err)
	fields = errx.ValidationErrors(err)
	be.Zero(t, fields)

	// Don't allocate for valid messages
	allocs := testing.AllocsPerRun(10, func() {
		var v errx.Validator
		v.AddIf(false, "field", "message: %d", 1)
		err = v.Err()
	})
	be.Equal(t, 0, allocs)
}
