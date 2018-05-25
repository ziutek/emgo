// +build f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package sdmmc

BUG: STM32F1 not supported.

In F1 SDIO is on AHB bus. Need to implement EnableClock, DisableClock and Reset (there is no AHBRST register).