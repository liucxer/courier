package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/liucxer/courier/codegen/formatx"
	"github.com/spf13/cobra"
)

func init() {
	cmdRoot.AddCommand(cmdFormat)
}

var cmdFormat = &cobra.Command{
	Use:   "fmt",
	Short: "format and sort imports",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()

		files := loadFiles(cwd, func(filename string) bool {
			return path.Ext(filename) == ".go" && !strings.Contains(filename, "/vendor/")
		})

		for _, filename := range files {
			fileInfo, _ := os.Stat(filename)
			bytes, _ := ioutil.ReadFile(filename)
			nextBytes, err := formatx.Format(filename, bytes, formatx.SortImportsProcess)
			if err != nil {
				panic(fmt.Errorf("errors %s in %s", filename, err.Error()))
			}
			if string(nextBytes) != string(bytes) {
				fmt.Printf("reformatted %s\n", filename)
				err := ioutil.WriteFile(filename, nextBytes, fileInfo.Mode())
				if err != nil {
					panic(err)
				}
			}
		}
	},
}

func loadFiles(dir string, filter func(filename string) bool) (filenames []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filename := path.Join(dir, file.Name())
		if file.IsDir() {
			filenames = append(filenames, loadFiles(filename, filter)...)
		} else {
			if filter(filename) {
				filenames = append(filenames, filename)
			}
		}
	}
	return
}
