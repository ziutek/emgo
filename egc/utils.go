package main

import (
	"os"
)

func logErr(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
}

func die(s string) {
	os.Stderr.WriteString(s + "\n")
	os.Exit(1)
}
