package handler

import (
	"context"
	"testing"

	"github.com/creachadair/jrpc2"
)

// Verify that the New function correctly handles the various type signatures
// it's advertised to support, and not others.
func TestNew(t *testing.T) {
	tests := []struct {
		v   interface{}
		bad bool
	}{
		{v: nil, bad: true},              // nil value
		{v: "not a function", bad: true}, // not a function

		// All the legal kinds...
		{v: func(context.Context) error { return nil }},
		{v: func(context.Context, *jrpc2.Request) (interface{}, error) { return nil, nil }},
		{v: func(context.Context) (int, error) { return 0, nil }},
		{v: func(context.Context, []int) error { return nil }},
		{v: func(context.Context, []bool) (float64, error) { return 0, nil }},
		{v: func(context.Context, ...int) error { return nil }},
		{v: func(context.Context, ...int) bool { return true }},
		{v: func(context.Context, ...string) (bool, error) { return false, nil }},
		{v: func(context.Context, *jrpc2.Request) error { return nil }},
		{v: func(context.Context, *jrpc2.Request) float64 { return 0 }},
		{v: func(context.Context, *jrpc2.Request) (byte, error) { return '0', nil }},
		{v: func(context.Context) bool { return true }},
		{v: func(context.Context, int) bool { return true }},

		// Things that aren't supposed to work.
		{v: func() error { return nil }, bad: true},                           // wrong # of params
		{v: func(a, b, c int) bool { return false }, bad: true},               // ...
		{v: func(byte) {}, bad: true},                                         // wrong # of results
		{v: func(byte) (int, bool, error) { return 0, true, nil }, bad: true}, // ...
		{v: func(string) error { return nil }, bad: true},                     // missing context
		{v: func(a, b string) error { return nil }, bad: true},                // P1 is not context
		{v: func(context.Context) (int, bool) { return 1, true }, bad: true},  // R2 is not error

		//lint:ignore ST1008 verify permuted error position does not match
		{v: func(context.Context) (error, float64) { return nil, 0 }, bad: true}, // ...
	}
	for _, test := range tests {
		got, err := newHandler(test.v)
		if !test.bad && err != nil {
			t.Errorf("newHandler(%T): unexpected error: %v", test.v, err)
		} else if test.bad && err == nil {
			t.Errorf("newHandler(%T): got %+v, want error", test.v, got)
		}
	}
}
