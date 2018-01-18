package compact

import (
	"github.com/thrift-iterator/go/protocol"
)

type TCompactType byte

const (
	TypeBooleanTrue  TCompactType = 0x01
	TypeBooleanFalse TCompactType = 0x02
	TypeByte         TCompactType = 0x03
	TypeI16          TCompactType = 0x04
	TypeI32          TCompactType = 0x05
	TypeI64          TCompactType = 0x06
	TypeDouble       TCompactType = 0x07
	TypeBinary       TCompactType = 0x08
	TypeList         TCompactType = 0x09
	TypeSet          TCompactType = 0x0A
	TypeMap          TCompactType = 0x0B
	TypeStruct       TCompactType = 0x0C
)

// TType value.
func (t TCompactType) ToTType() protocol.TType {
	switch TCompactType(byte(t) & 0x0f) {
	case TypeBooleanFalse, TypeBooleanTrue:
		return protocol.BOOL
	case TypeByte:
		return protocol.BYTE
	case TypeI16:
		return protocol.I16
	case TypeI32:
		return protocol.I32
	case TypeI64:
		return protocol.I64
	case TypeDouble:
		return protocol.DOUBLE
	case TypeBinary:
		return protocol.STRING
	case TypeList:
		return protocol.LIST
	case TypeSet:
		return protocol.SET
	case TypeMap:
		return protocol.MAP
	case TypeStruct:
		return protocol.STRUCT
	}
	return protocol.STOP
}
