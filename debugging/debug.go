package debugging

import (
	"io"
	"io/ioutil"
	"os"
)

var DebugOut io.Writer = ioutil.Discard

func EnableDebug() {
	DebugOut = os.Stderr
}
