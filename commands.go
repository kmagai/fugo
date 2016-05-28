package main

import (
	"github.com/kmagai/fugo/subcmd"
	"github.com/mitchellh/cli"
)

func commands(style *subcmd.Style) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"check": func() (cli.Command, error) {
			return &subcmd.Check{
				Style: *style,
			}, nil
		},
		"add": func() (cli.Command, error) {
			return &subcmd.Add{
				Style: *style,
			}, nil
		},
		"remove": func() (cli.Command, error) {
			return &subcmd.Remove{
				Style: *style,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &subcmd.VersionCommand{
				Style:    *style,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
