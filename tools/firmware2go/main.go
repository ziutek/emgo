package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 || os.Args[1] != "bytes" && os.Args[1] != "string" {
		os.Stderr.WriteString("Usage: firmware2go {bytes|string} BINARY_FILE\n")
		os.Exit(1)
	}
	data, err := ioutil.ReadFile(os.Args[2])
	checkErr(err)
	w := bufio.NewWriter(os.Stdout)
	_, err = w.WriteString("package main\n\n")
	checkErr(err)
	if os.Args[1] == "bytes" {
		_, err = w.WriteString("//emgo:const\nvar firmware = [...]byte{")
		checkErr(err)
		for i, b := range data {
			if i%15 == 0 {
				_, err = w.WriteString("\n\t")
				checkErr(err)
			}
			_, err = fmt.Fprintf(w, " %d,", b)
			checkErr(err)
		}
		_, err = w.WriteString("\n}\n")
		checkErr(err)
	} else {
		_, err = w.WriteString("const firmware = \"")
		checkErr(err)
		for i, b := range data {
			if i%18 == 0 {
				_, err = w.WriteString("\" +\n\t\"")
				checkErr(err)
			}
			_, err = fmt.Fprintf(w, "\\x%02X", b)
			checkErr(err)
		}
		_, err = w.WriteString("\"\n")
		checkErr(err)
	}
	checkErr(w.Flush())
}
