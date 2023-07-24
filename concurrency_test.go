package jrpc2_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/server"
)

func TestConcurrency(t *testing.T) {

	t.Run("correctly propagates context cancel", func(t *testing.T) {

		done := make(chan bool)

		loc := server.NewLocal(handler.Map{
			"test": handler.New(func(ctx context.Context) error {
				<-ctx.Done()
				<-done
				return nil
			}),
		}, &server.LocalOptions{
			Server: &jrpc2.ServerOptions{Concurrency: 1},
		})

		ctx, cancel := context.WithCancel(context.Background())
		go loc.Client.Call(ctx, "test", nil)
		cancel()

		select {
		case <-done:
			//ok
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout!")
		}
	})

	t.Run("doesn't get clogged with unfinished calls", func(t *testing.T) {

		started := false

		loc := server.NewLocal(handler.Map{
			"test1": handler.New(func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			}),
			"test2": handler.New(func(ctx context.Context) error {
				started = true
				return nil
			}),
		}, &server.LocalOptions{
			Server: &jrpc2.ServerOptions{Concurrency: 1},
		})

		ctx, cancel := context.WithCancel(context.Background())
		go loc.Client.Call(ctx, "test1", nil)
		runtime.Gosched()
		cancel()

		done := make(chan bool)

		go func() {
			loc.Client.Call(context.Background(), "test2", nil)
			done <- true
		}()

		select {
		case <-done:
			//ok
		case <-time.After(time.Millisecond * 100):
		}

		if !started {
			t.Fatal("expected the second handler to have been started")
		}
	})
}
