package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/cadyyan/void-tool/cmd"
)

func main() {
	ctx := context.Background()

	ctx, cancelSignals := signal.NotifyContext(ctx, os.Interrupt)
	defer cancelSignals()

	rootCommand := cmd.NewRootCommand()
	if err := rootCommand.ExecuteContext(ctx); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
