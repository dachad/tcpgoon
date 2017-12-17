package main

import (
	"github.com/dachad/tcpgoon/cmd"
)

var buildstamp = "No build time provided"
var githash = "No git hash provided"

func main() {
	cmd.Execute(buildstamp, githash)
}
