package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func escape(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return r
	}
	return '_'
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printf(w *bufio.Writer, f string, a ...interface{}) {
	_, err := fmt.Fprintf(w, f, a...)
	checkErr(err)
}

func flush(w *bufio.Writer) {
	checkErr(w.Flush())
}
