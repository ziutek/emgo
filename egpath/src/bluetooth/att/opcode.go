package att

type Method byte

const (
	ErrorRsp       Method = 0x01
	ReqdByTypeReq  Method = 0x08
	ReadByTypeRsp  Method = 0x09
	ReadByGroupReq Method = 0x10
	ReadByGroupRsp Method = 0x11
)
