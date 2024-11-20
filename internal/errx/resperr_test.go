package errx_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/carlmjohnson/be"
	"github.com/spotlightpa/moreofa/internal/errx"
)

func TestGetCode(t *testing.T) {
	base := errx.WithStatusCode(errors.New(""), 5)
	wrapped := fmt.Errorf("wrapping: %w", base)

	testCases := map[string]struct {
		error
		int
	}{
		"nil":         {nil, 200},
		"default":     {errors.New(""), 500},
		"set":         {errx.WithStatusCode(errors.New(""), 3), 3},
		"set-nil":     {errx.WithStatusCode(nil, 4), 4},
		"wrapped":     {wrapped, 5},
		"set-message": {errx.WithUserMessage(nil, "xxx"), 400},
		"set-both":    {errx.WithCodeAndMessage(nil, 6, "xx"), 6},
		"context":     {context.DeadlineExceeded, 504},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			be.Equal(t, tc.int, errx.StatusCode(tc.error))
		})
	}
}

func TestSetCode(t *testing.T) {
	t.Run("same-message", func(t *testing.T) {
		err := errors.New("hello")
		coder := errx.WithStatusCode(err, 400)
		got := coder.Error()
		want := "[400] " + err.Error()
		be.Equal(t, want, got)
	})
	t.Run("keep-chain", func(t *testing.T) {
		err := errors.New("hello")
		coder := errx.WithStatusCode(err, 3)
		be.True(t, errors.Is(coder, err))
	})
	t.Run("set-nil", func(t *testing.T) {
		coder := errx.WithStatusCode(nil, 400)
		be.In(t, http.StatusText(400), coder.Error())
	})
	t.Run("override-default", func(t *testing.T) {
		err := context.DeadlineExceeded
		coder := errx.WithStatusCode(err, 3)
		code := errx.StatusCode(coder)
		be.Equal(t, 3, code)
	})
}

func TestGetMsg(t *testing.T) {
	base := errx.WithUserMessage(errors.New(""), "5")
	wrapped := fmt.Errorf("wrapping: %w", base)

	testCases := map[string]struct {
		error
		string
	}{
		"nil":     {nil, ""},
		"default": {errors.New(""), "Internal Server Error"},
		"set":     {errx.WithUserMessage(errors.New(""), "3"), "3"},
		"set-nil": {errx.WithUserMessage(nil, "4"), "4"},
		"wrapped": {wrapped, "5"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			be.Equal(t, tc.string, errx.UserMessage(tc.error))
		})
	}
}

func TestSetMsg(t *testing.T) {
	t.Run("same-message", func(t *testing.T) {
		err := errors.New("hello")
		msgr := errx.WithUserMessage(err, "a")
		be.Equal(t, err.Error(), msgr.Error())
	})
	t.Run("keep-chain", func(t *testing.T) {
		err := errors.New("hello")
		msgr := errx.WithUserMessage(err, "a")
		be.True(t, errors.Is(msgr, err))
	})
	t.Run("set-nil", func(t *testing.T) {
		msgr := errx.WithUserMessage(nil, "a")
		be.Equal(t, "UserMessage<a>", msgr.Error())
	})
}

func TestMsgf(t *testing.T) {
	msg := "hello 1, 2, 3"
	err := errx.WithUserMessagef(nil, "hello %d, %d, %d", 1, 2, 3)
	be.Equal(t, msg, errx.UserMessage(err))
}

func TestNotFound(t *testing.T) {
	path := "/example/url"
	r, _ := http.NewRequest(http.MethodGet, path, nil)
	err := errx.NotFound(r)
	be.In(t, path, err.Error())
	be.In(t, path, errx.UserMessage(err))
	be.Equal(t, 404, errx.StatusCode(err))
}

func TestNew(t *testing.T) {
	t.Run("flat", func(t *testing.T) {
		err := errx.New(404, "hello %s", "world")
		be.Equal(t, "Not Found", errx.UserMessage(err))
		be.Equal(t, 404, errx.StatusCode(err))
		be.Equal(t, "[404] hello world", err.Error())
	})
	t.Run("chain", func(t *testing.T) {
		const setMsg = "msg1"
		inner := errx.WithUserMessage(nil, setMsg)
		w1 := errx.New(5, "w1: %w", inner)
		w2 := errx.New(6, "w2: %w", w1)
		be.Equal(t, setMsg, errx.UserMessage(w2))
		be.Equal(t, 5, errx.StatusCode(w1))
		be.Equal(t, 6, errx.StatusCode(w2))
		be.Equal(t, "[6] w2: [5] w1: UserMessage<msg1>", w2.Error())
	})
}
