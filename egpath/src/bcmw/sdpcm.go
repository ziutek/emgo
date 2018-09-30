package bcmw

import (
	"encoding/binary/le"
	"sdcard"
)

// SDPCM is Broadcom protocol that allows to multiplex WLAN data frames, I/O
// controll function and asynchronous event signaling.

const (
	sdpcmGet = 0
	sdpcmSet = 2
)

const (
	wlcGetMagic                   = 0
	wlcGetVersion                 = 1
	wlcUp                         = 2
	wlcDown                       = 3
	wlcGetLoop                    = 4
	wlcSetLoop                    = 5
	wlcDump                       = 6
	wlcGetMsgLevel                = 7
	wlcSetMsgLevel                = 8
	wlcGetPromisc                 = 9
	wlcSetPromisc                 = 10
	wlcGetRate                    = 12
	wlcGetInstance                = 14
	wlcGetInfra                   = 19
	wlcSetInfra                   = 20
	wlcGetAuth                    = 21
	wlcSetAuth                    = 22
	wlcGetBSSID                   = 23
	wlcSetBSSID                   = 24
	wlcGetSSID                    = 25
	wlcSetSSID                    = 26
	wlcRestart                    = 27
	wlcGetChannel                 = 29
	wlcSetChannel                 = 30
	wlcGetSRL                     = 31
	wlcSetSRL                     = 32
	wlcGetLRL                     = 33
	wlcSetLRL                     = 34
	wlcGetPLCPHdr                 = 35
	wlcSetPLCPHdr                 = 36
	wlcGetRadio                   = 37
	wlcSetRadio                   = 38
	wlcGetPhyType                 = 39
	wlcDumpRate                   = 40
	wlcSetRateParams              = 41
	wlcGetKey                     = 44
	wlcSetKey                     = 45
	wlcGetRegulatory              = 46
	wlcSetRegulatory              = 47
	wlcGetPasiveScan              = 48
	wlcSetPasiveScan              = 49
	wlcScan                       = 50
	wlcScanResults                = 51
	wlcDisassoc                   = 52
	wlcReassoc                    = 53
	wlcGetRoamTrigger             = 54
	wlcSetRoamTrigger             = 55
	wlcGetRoamDelta               = 56
	wlcSetRoamDelta               = 57
	wlcGetRoamScanPeriod          = 58
	wlcSetRoamScanPeriod          = 59
	wlcEVM                        = 60
	wlcGetTxAnt                   = 61
	wlcSetTxAnt                   = 62
	wlcGetAntDiv                  = 63
	wlcSetAntDiv                  = 64
	wlcGetClosed                  = 67
	wlcSetClosed                  = 68
	wlcGetMACList                 = 69
	wlcSetMACList                 = 70
	wlcGetRateSet                 = 71
	wlcSetRateSet                 = 72
	wlcLongTrain                  = 74
	wlcGetBeaconPeriod            = 75
	wlcSetBeaconPeriod            = 76
	wlcGetDTIMPeriod              = 77
	wlcSetDTIMPeriod              = 78
	wlcGetSROM                    = 79
	wlcSetSROM                    = 80
	wlcGetWEPRestrict             = 81
	wlcSetWEPRestrict             = 82
	wlcGetCountry                 = 83
	wlcSetCountry                 = 84
	wlcGetPM                      = 85
	wlcSetPM                      = 86
	wlcGetWake                    = 87
	wlcSetWake                    = 88
	wlcGetForceLink               = 90
	wlcSetForceLink               = 91
	wlcFreqAccuracy               = 92
	wlcCarrierSuppress            = 93
	wlcGetPhyReg                  = 94
	wlcSetPhyReg                  = 95
	wlcGetRadioReg                = 96
	wlcSetRadioReg                = 97
	wlcGetRevInfo                 = 98
	wlcGetUCAntDiv                = 99
	wlcSetUCAntDiv                = 100
	wlcRReg                       = 101
	wlcWReg                       = 102
	wlcGetMACMode                 = 105
	wlcSetMACMode                 = 106
	wlcGetMonitor                 = 107
	wlcSetMonitor                 = 108
	wlcGetGMode                   = 109
	wlcSetGMode                   = 110
	wlcGetLegacyERP               = 111
	wlcSetLegacyERP               = 112
	wlcGetRxAnt                   = 113
	wlcGetCurrRateSet             = 114
	wlcGetScanSuppress            = 115
	wlcSetScanSuppress            = 116
	wlcGetAP                      = 117
	wlcSetAP                      = 118
	wlcGetEAPRestrict             = 119
	wlcSetEAPRestrict             = 120
	wlcSCBAuthorize               = 121
	wlcSCBDeauthorize             = 122
	wlcGetWDSList                 = 123
	wlcSetWDSList                 = 124
	wlcGetATIM                    = 125
	wlcSetATIM                    = 126
	wlcGetRSSI                    = 127
	wlcGetPhyAntDiv               = 128
	wlcSetPhyAntDiv               = 129
	wlcAPRxOnly                   = 130
	wlcGetTxPathPwr               = 131
	wlcSetTxPathPwr               = 132
	wlcGetWSEC                    = 133
	wlcSetWSEC                    = 134
	wlcGetPhyNoise                = 135
	wlcGetBSSInfo                 = 136
	wlcGetPktCnts                 = 137
	wlcGetLazyWDS                 = 138
	wlcSetLazyWDS                 = 139
	wlcGetBandList                = 140
	wlcGetBand                    = 141
	wlcSetBand                    = 142
	wlcSCBDeauthenticate          = 143
	wlcGetShortSlot               = 144
	wlcGetShortSlotOverride       = 145
	wlcSetShortSlotOverride       = 146
	wlcGetShortSlotRestrict       = 147
	wlcSetShortSlotRestrict       = 148
	wlcGetGModeProtection         = 149
	wlcGetGModeProtectionOverride = 150
	wlcSetGModeProtectionOverride = 151
	wlcUpgrade                    = 152
	wlcGetIgnoreBCNS              = 155
	wlcSetIgnoreBCNS              = 156
	wlcGetSCBTimeout              = 157
	wlcSetSCBTimeout              = 158
	wlcGetAssocList               = 159
	wlcGetClk                     = 160
	wlcSetClk                     = 161
	wlcGetUp                      = 162
	wlcOut                        = 163
	wlcGetWPAAuth                 = 164
	wlcSetWPAAuth                 = 165
	wlcGetUCFlags                 = 166
	wlcSetUCFlags                 = 167
	wlcGetPwrIdx                  = 168
	wlcSetPwrIdx                  = 169
	wlcGetTSSI                    = 170
	wlcGetSupRatesetOverride      = 171
	wlcSetSupRatesetOverride      = 172
	wlcGetProtectionControl       = 178
	wlcSetProtectionControl       = 179
	wlcGetPhyList                 = 180
	wlcEncryptStrength            = 181
	wlcDecyptStatus               = 182
	wlcGetKeySeq                  = 183
	wlcGetScanChannelTime         = 184
	wlcSetScanChannelTime         = 185
	wlcGetScanUnassocTime         = 186
	wlcSetScanUnassocTime         = 187
	wlcGetScanHomeTime            = 188
	wlcSetScanHomeTime            = 189
	wlcGetScanNProbes             = 190
	wlcSetScanNProbes             = 191
	wlcGetPrbRespTimeout          = 192
	wlcSetPrbRespTimeout          = 193
	wlcGetAtten                   = 194
	wlcSetAtten                   = 195
	wlcGetSHMem                   = 196
	wlcSetSHMem                   = 197
	wlcSetWsecTest                = 200
	wlcSCBDeauthenticateForReason = 201
	wlcTKIPCounterMeasures        = 202
	wlcGetPIOMode                 = 203
	wlcSetPIOMode                 = 204
	wlcSetAssocPrefer             = 205
	wlcGetAssocPrefer             = 206
	wlcSetRoamPrefer              = 207
	wlcGetRoamPrefer              = 208
	wlcSetLED                     = 209
	wlcGetLED                     = 210
	wlcGetInterferenceMode        = 211
	wlcSetInterferenceMode        = 212
	wlcGetChannelQA               = 213
	wlcStartChannelQA             = 214
	wlcGetChannelSel              = 215
	wlcStartChannelSel            = 216
	wlcGetValidChannels           = 217
	wlcGetFakeFrag                = 218
	wlcSetFakeFrag                = 219
	wlcGetPwrOutPercentage        = 220
	wlcSetPwrOutPercentage        = 221
	wlcSetBadFramePreempt         = 222
	wlcGetBadFramePreempt         = 223
	wlcSetLeapList                = 224
	wlcGetLeapList                = 225
	wlcGetCWMin                   = 226
	wlcSetCWMin                   = 227
	wlcGetCWMax                   = 228
	wlcSetCWMax                   = 229
	wlcGetWET                     = 230
	wlcSetWET                     = 231
	wlcGetPub                     = 232
	wlcGetKeyPrimary              = 235
	wlcSetKeyPrimary              = 236
	wlcGetACIArgs                 = 238
	wlcSetACIArgs                 = 239
	wlcUnsetCallback              = 240
	wlcSetCallback                = 241
	wlcGetRadar                   = 242
	wlcSetRadar                   = 243
	wlcSetSpectManagement         = 244
	wlcGetSpectManagement         = 245
	wlcWDSGetRemoteHwAddr         = 246
	wlcWDSGetWPASup               = 247
	wlcSetCSScanTimer             = 248
	wlcGetCSScanTimer             = 249
	wlcMeasureRequest             = 250
	wlcInit                       = 251
	wlcSendQuiet                  = 252
	wlcKeepalive                  = 253
	wlcSrndPweConstraint          = 254
	wlcUpgradeStatus              = 255
	wlc_CURRENT_PWR               = 256
	wlcGetScanPassiveTime         = 257
	wlcSetScanPassiveTime         = 258
	wlcLegacyLinkBehavior         = 259
	wlcGetChannelInCountry        = 260
	wlcGetCountryList             = 261
	wlcGetVar                     = 262
	wlcSetVar                     = 263
	wlcNVRAMGet                   = 264
	wlcNVRAMSet                   = 265
	wlcNVRAMDump                  = 266
	wlcReboot                     = 267
	wlcSetWSecPMK                 = 268
	wlcGetAuthMode                = 269
	wlcSetAuthMode                = 270
	wlcGetWakeEntry               = 271
	wlcSetWakeEntry               = 272
	wlcNDConfigItem               = 273
	wlcNVOTPW                     = 274
	wlcOTPW                       = 275
	wlcIOVBlockGet                = 276
	wlcIOVModulesGet              = 277
	wlcSoftReset                  = 278
	wlcGetAllowMode               = 279
	wlcSetAllowMode               = 280
	wlcGetDesiredBSSID            = 281
	wlcSetDesiredBSSID            = 282
	wlc_DISASSOC_MYAP             = 283
	wlcGetNBands                  = 284
	wlcGetBandStates              = 285
	wlcGetWLCBSSInfo              = 286
	wlcGetAssocInfo               = 287
	wlcGetOIDPhy                  = 288
	wlcSetOIDPhy                  = 289
	wlcSetAssocTime               = 290
	wlcGetDesiredSSID             = 291
	wlcGetChanSpec                = 292
	wlcGetAssocState              = 293
	wlcSetPhyState                = 294
	wlcGetScanPending             = 295
	wlcGetScanReqPending          = 296
	wlcGetPrevRoamReason          = 297
	wlcSetPrevRoamReason          = 298
	wlcGetBandStatesPI            = 299
	wlcGetPhyState                = 300
	wlcGetBssWPARSN               = 301
	wlcGetBssWPA2RSN              = 302
	wlcGetBssBCNTS                = 303
	wlcGetIntDosassoc             = 304
	wlcSetNumPeers                = 305
	wlcGetNumBSS                  = 306
	wlcGetWSecPMK                 = 318
	wlcGetRandomBytes             = 319
	wlcLast                       = 320
)

const (
	ifSta = 0
	ifAP  = 1
	ifP2P = 2
)

func (d *Driver) sdpcmReadFrame() bool {
	if d.error() {
		return false
	}
	sd := d.sd
	var buf [1]uint64
	sd.SetupData(sdcard.Recv|sdcard.IO|sdcard.Block4, buf[:], 4)
	_, d.ioStatus = sd.SendCmd(sdcard.CMD53(wlanData, 0, sdcard.Read, 4)).R5()
	if d.error() {
		return false
	}
	bytes := sdcard.Data(buf[:]).Bytes()
	size := le.Decode16(bytes)
	checksum := le.Decode16(bytes[2:])
	if size|checksum == 0 {
		d.debug("no more frames")
		return false
	}
	d.debug("size: %d (0x%04X) checksum: 0x%04X\n", size, size, checksum)
	if size != ^checksum {
		d.debug("bad header checksum: 0x%04X != 0x%04X\n", size, ^checksum)
		return false
	}
	return true
}
