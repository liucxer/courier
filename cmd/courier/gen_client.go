package main

import (
	"net/url"
	"os"
	"path"

	"github.com/liucxer/courier/httptransport/client/generator"
	"github.com/liucxer/courier/packagesx"
	"github.com/spf13/cobra"
)

var (
	cmdGenClientFlagFile    string
	cmdGenClientFlagSpecURL string
)

func init() {
	cmdGen.AddCommand(cmdGenClient)

	cmdGenClient.Flags().
		StringVarP(&cmdGenClientFlagSpecURL, "spec-url", "", "", "client spec url")
	cmdGenClient.Flags().
		StringVarP(&cmdGenClientFlagFile, "file", "", "", "client spec file")

}

var cmdGenClient = &cobra.Command{
	Use:     "client",
	Example: "client demo",
	Short:   "generate client by open api",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			return
		}
		u := &url.URL{}

		if cmdGenClientFlagFile != "" {
			u.Scheme = "file"
			cwd, _ := os.Getwd()
			u.Path = path.Join(cwd, cmdGenClientFlagFile)
		}

		if cmdGenClientFlagSpecURL != "" {
			uri, err := url.Parse(cmdGenClientFlagSpecURL)
			if err != nil {
				panic(err)
			}
			u = uri
		}

		runGenerator(func(pkg *packagesx.Package) Generator {
			g := generator.NewClientGenerator(args[0], u)
			g.Load()
			return g
		})
	},
}
