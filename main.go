package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"toolman.org/base/toolman/v2"
)

var program, logdir string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic("Cannot determine current user: " + err.Error())
	}

	program = filepath.Base(os.Args[0])
	logdir = filepath.Join(usr.HomeDir, ".log", "polybar", program)
}

func main() {
	cfg := defaults()
	toolman.Init(cfg.flags(), toolman.LogFlushInterval(500*time.Millisecond), toolman.LogDir(logdir), toolman.MakeLogDir())

	ctx := context.Background()

	if err := run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg *configuration) error {
	oh := &openhab{cfg.ohBase}

	cs, err := oh.currentState(ctx)
	if err != nil {
		return err
	}

	som := cfg.stateOutputMap()

	fmt.Println(som[cs])

	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ch := make(chan fanState)

	eg.Go(func() error { return oh.events(ctx, ch) })

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGUSR1)

	eg.Go(func() error {
		for {
			select {
			case fs := <-ch:
				fmt.Println(som[fs])
			case <-sig:
				fmt.Println(som[fsChange])
				if err := oh.toggleState(ctx); err != nil {
					return err
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	return eg.Wait()
}
