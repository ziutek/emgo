package main

import (
	"os"
)

func logErr(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
}
