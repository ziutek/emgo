// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package internal

import (
	"unsafe"
)

func APB_SetLPEnabled(_ unsafe.Pointer, _ bool) {}
