package main

// Copyright (c) 2015 Broadcom
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// thislist of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// 3. Neither the name of Broadcom nor the names of other contributors to this
// software may be used to endorse or promote products derived from this
// software without specific prior written permission.
//
// 4. This software may not be used as a standalone product, and may only be
// used as incorporated in your product or device that incorporates Broadcom
// wireless connectivity products and solely for the purpose of enabling the
// functionalities of such Broadcom products.
//
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY WARRANTIES OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NON-INFRINGEMENT, ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
// OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
// EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// BCM43362 NVRAM variables for WM-N-BM-09 USI SiP

const nvram = "" +
	"manfid=0x2d0" + "\x00" +
	"prodid=0x492" + "\x00" +
	"vendid=0x14e4" + "\x00" +
	"devid=0x4343" + "\x00" +
	"boardtype=0x0636" + "\x00" +
	"boardrev=0x1201" + "\x00" +
	"boardnum=777" + "\x00" +
	"xtalfreq=26000" + "\x00" +
	"boardflags=0xa00" + "\x00" +
	"sromrev=3" + "\x00" +
	"wl0id=0x431b" + "\x00" +
	"macaddr=02:0A:F7:e0:ae:ce" + "\x00" +
	"aa2g=3" + "\x00" +
	"ag0=2" + "\x00" +
	"maxp2ga0=74" + "\x00" +
	"ofdm2gpo=0x44111111" + "\x00" +
	"mcs2gpo0=0x4444" + "\x00" +
	"mcs2gpo1=0x6444" + "\x00" +
	"pa0maxpwr=80" + "\x00" +
	"pa0b0=5264" + "\x00" +
	"pa0b1=64897" + "\x00" +
	"pa0b2=65359" + "\x00" +
	"pa0itssit=62" + "\x00" +
	"pa1itssit=62" + "\x00" +
	"temp_based_dutycy_en=1" + "\x00" +
	"tx_duty_cycle_ofdm=100" + "\x00" +
	"tx_duty_cycle_cck=100" + "\x00" +
	"tx_ofdm_temp_0=115" + "\x00" +
	"tx_cck_temp_0=115" + "\x00" +
	"tx_ofdm_dutycy_0=40" + "\x00" +
	"tx_cck_dutycy_0=40" + "\x00" +
	"tx_ofdm_temp_1=255" + "\x00" +
	"tx_cck_temp_1=255" + "\x00" +
	"tx_ofdm_dutycy_1=40" + "\x00" +
	"tx_cck_dutycy_1=40" + "\x00" +
	"tx_tone_power_index=40" + "\x00" +
	"tx_tone_power_index.fab.3=48" + "\x00" +
	"cckPwrOffset=0" + "\x00" +
	"ccode=0" + "\x00" +
	"rssismf2g=0xa" + "\x00" +
	"rssismc2g=0x3" + "\x00" +
	"rssisav2g=0x7" + "\x00" +
	"triso2g=0" + "\x00" +
	"noise_cal_enable_2g=0" + "\x00" +
	"noise_cal_po_2g=0" + "\x00" +
	"noise_cal_po_2g.fab.3=-2" + "\x00" +
	"swctrlmap_2g=0x0a030a03,0x0c050c05,0x0c050c05,0x0,0x1ff" + "\x00" +
	"temp_add=29767" + "\x00" +
	"temp_mult=425" + "\x00" +
	"temp_q=10" + "\x00" +
	"initxidx2g=45" + "\x00" +
	"tssitime=1" + "\x00" +
	"rfreg033=0x19" + "\x00" +
	"rfreg033_cck=0x1f" + "\x00" +
	"cckPwrIdxCorr=-8" + "\x00" +
	"spuravoid_enable2g=1" + "\x00" +
	"edonthd=-70" + "\x00" +
	"edoffthd=-76" + "\x00" +
	"\x00\x00"
