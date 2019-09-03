package main

import (
	"fmt"
	"github.com/adrianrudnik/n26-mt940-converter/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
