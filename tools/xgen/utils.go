package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func fdie(f, format string, args ...interface{}) {
	die(f+": "+format, args...)
}

func checkErr(err error) {
	if err != nil {
		die("Error: %v", err)
	}
}

func nameval(line string, sep byte) (name, val string) {
	i := strings.IndexByte(line, sep)
	if i < 0 {
		return
	}
	name = strings.TrimSpace(line[:i])
	line = strings.TrimSpace(line[i+1:])
	i = strings.IndexFunc(line, unicode.IsSpace)
	if i < 0 {
		val = line
	} else {
		val = line[:i]
	}
	return
}
