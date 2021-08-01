package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/liucxer/courier/packagesx"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use: "courier",
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func logCost() func() {
	startedAt := time.Now()

	return func() {
		log.Printf("costs %s", color.GreenString("%0.0f ms", float64(time.Now().Sub(startedAt)/time.Millisecond)))
	}
}

type Generator interface {
	Output(cwd string)
}

func runGenerator(createGenerator func(pkg *packagesx.Package) Generator) {
	defer logCost()()

	cwd, _ := os.Getwd()
	pkg, err := packagesx.Load(cwd)
	if err != nil {
		panic(err)
	}

	g := createGenerator(pkg)
	g.Output(cwd)
}
