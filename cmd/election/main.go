package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"

	"github.com/central-university-dev/2024-spring-go-course-lesson8-leader-election/internal/commands"
)

func main() {
	eg, _ := errgroup.WithContext(context.Background())
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
