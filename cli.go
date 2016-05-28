package main

import (
	"fmt"
	"os"

	"github.com/kmagai/fugo/subcmd"
	"github.com/mitchellh/cli"
)

func run(args []string) int {
	style := &subcmd.Style{
		UI: &cli.ColoredUi{
			InfoColor:  cli.UiColorBlue,
			ErrorColor: cli.UiColorRed,
			Ui: &cli.BasicUi{
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
				Reader:      os.Stdin,
			},
		}}

	return runCustom(args, commands(style))
}

func runCustom(args []string, comm map[string]cli.CommandFactory) int {
	// treat version as if it were a subcmd
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	// default action to check
	if len(args) == 0 {
		checkArg := make([]string, 1)
		checkArg[0] = "check"
		args = checkArg
	}

	cli := &cli.CLI{
		Args:       args,
		Commands:   comm,
		Version:    Version,
		HelpFunc:   cli.BasicHelpFunc(Name),
		HelpWriter: os.Stdout,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed to execute: %s\n", err.Error())
	}

	return exitCode
}
