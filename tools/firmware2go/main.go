package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	if len(os.Args) != 4 || (os.Args[1] != "le" && os.Args[1] != "be") {
		fmt.Fprintf(
			os.Stderr, "Usage: firmware2go le|be INPUT.bin OUTPUT.go\n",
		)
		os.Exit(1)
	}
	le := os.Args[1] != "le"
	f, err := os.Open(os.Args[2])
	checkErr(err)
	inp := bufio.NewReader(f)
	f, err = os.Create(os.Args[3])
	checkErr(err)
	defer f.Close()
	out := bufio.NewWriter(f)
	_, err = out.WriteString(
		"package main\n" + "//emgo:const\n" + "var firmware = [...]uint64{\n",
	)
	checkErr(err)
	var buf [8]byte
	for i := 0; ; i++ {
		n, err := io.ReadFull(inp, buf[:])
		if err != nil && err != io.ErrUnexpectedEOF {
			if err == io.EOF {
				if i%3 != 0 {
					checkErr(out.WriteByte('\n'))
				}
				break
			}
			checkErr(err)
		}
		var v uint64
		if le {
			for k := 0; k < n; k++ {
				v |= uint64(buf[k]) << uint(8*k)
			}
		} else {
			for k := 0; k < n; k++ {
				v |= uint64(buf[k]) << uint(56-8*k)
			}
		}
		_, err = fmt.Fprintf(out, " 0x%016X,", v)
		checkErr(err)
		if i%3 == 2 {
			checkErr(out.WriteByte('\n'))
		}
	}
	_, err = out.WriteString("}\n")
	checkErr(err)
	checkErr(out.Flush())
}
