package main

import (
	"fmt"

	"github.com/dachad/tcpgoon/cmd"
)

var version string

func main() {
	fmt.Println("Running tcpgoon version " + version)
	cmd.Execute()
}
