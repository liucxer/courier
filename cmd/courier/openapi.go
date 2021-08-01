package main

import (
	"github.com/liucxer/courier/httptransport/openapi/generator"
	"github.com/liucxer/courier/packagesx"
	"github.com/spf13/cobra"
)

var cmdSwagger = &cobra.Command{
	Use:     "openapi",
	Aliases: []string{"swagger"},
	Short:   "scan current project and generate openapi.json",
	Run: func(cmd *cobra.Command, args []string) {
		runGenerator(func(pkg *packagesx.Package) Generator {
			g := generator.NewOpenAPIGenerator(pkg)
			g.Scan()
			return g
		})
	},
}

func init() {
	cmdRoot.AddCommand(cmdSwagger)
}
