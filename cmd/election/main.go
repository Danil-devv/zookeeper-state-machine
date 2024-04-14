package main

import (
	"context"
	"fmt"
	"hw/internal/commands"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	eg := errgroup.Group{}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan os.Signal, 1)
	signal.Notify(errCh, os.Interrupt)
	go func() {
		<-errCh
		fmt.Println("got SIGINT notify")
		cancel()
	}()

	rootCmd, err := commands.InitRunCommand(ctx)
	if err != nil {
		fmt.Println("init run command: %w", err)
		os.Exit(1)
	}

	eg.Go(func() error {
		return rootCmd.Execute()
	})

	if err = eg.Wait(); err != nil {
		fmt.Println("run command: %w", err)
		os.Exit(1)
	}
}
