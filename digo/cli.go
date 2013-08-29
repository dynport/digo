package main

import (
	"fmt"
	"github.com/dynport/digo"
	"github.com/dynport/gocli"
)

func init() {
	cli.Register("version",
		&gocli.Action {
			Description: "Print version and revision",
			Handler: VersionAction,
		},
		)
}

func VersionAction(args *gocli.Args) error {
	table := gocli.NewTable()
	table.Add("Version", digo.VERSION)
	table.Add("Revision", digo.GITCOMMIT)
	fmt.Printf(table.String())
	return nil
}
