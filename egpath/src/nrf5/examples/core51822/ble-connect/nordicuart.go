package main

import (
	"bluetooth/att"
	"bluetooth/uuid"
)

//emgo:const
var (
	serviceNordicUART = uuid.Long{0x6E400001B5A3F393, 0xE0A9E50E24DCCA9E}
)

type nordicUART struct {
	_ byte
}

func (u *nordicUART) ServeATT(w *att.ResponseWriter, r *att.Request) {

}
