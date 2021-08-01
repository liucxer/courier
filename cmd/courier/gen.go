package main

import (
	"github.com/spf13/cobra"
)

func init() {
	cmdRoot.AddCommand(cmdGen)
}

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "generators",
}
