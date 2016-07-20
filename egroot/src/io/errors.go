package io

import "errors"

var (
	EOF              = errors.New("EOF")
	ErrUnexpectedEOF = errors.New("unexpected EOF")
	ErrShortWrite    = errors.New("short write")
)
