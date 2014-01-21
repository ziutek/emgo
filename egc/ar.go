package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

var arHeader = []byte("!<arch>\n")

func arReadFile(apath string, name string) ([]byte, error) {
	a, err := os.Open(apath)
	if err != nil {
		return nil, err
	}
	defer a.Close()

	blen := 16 + 12 + 6 + 6 + 8 + 10 + 2
	buf := make([]byte, blen)

	n, err := io.ReadFull(a, buf[:len(arHeader)])

	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	if err == io.ErrUnexpectedEOF || !bytes.Equal(buf[:n], arHeader) {
		err = fmt.Errorf(
			"%s is too short or doesn't begin from ar header", apath,
		)
		return nil, err
	}

	bname := []byte(name)

	for {
		n, err = io.ReadFull(a, buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				err = fmt.Errorf(
					"archive %s doesn't contain %s file",
					apath, bname,
				)
			}
			return nil, err
		}
		if buf[blen-2] != '`' || buf[blen-1] != '\n' {
			err = fmt.Errorf("bad file header magic in %s", apath)
			return nil, err
		}
		fname := bytes.TrimRight(buf[:16], " ")
		if last := len(fname) - 1; fname[last] == '/' {
			// GNU ar
			fname = fname[:last]
		}
		flen, err := strconv.ParseUint(
			string(bytes.TrimRight(buf[48:58], " ")),
			10, 64,
		)
		if err != nil {
			err = fmt.Errorf(
				"bad file size for %s in %s: %s",
				fname, apath, err,
			)
			return nil, err
		}

		if bytes.Equal(fname, bname) {
			buf = make([]byte, flen)
			if _, err = io.ReadFull(a, buf); err != nil {
				err = fmt.Errorf(
					"can't read %s file from %s: %s",
					fname, apath, err,
				)
				return nil, err
			}
			return buf, nil
		}

		if flen&1 != 0 {
			flen++
		}
		if _, err = a.Seek(int64(flen), 1); err != nil {
			return nil, err
		}
	}
}
