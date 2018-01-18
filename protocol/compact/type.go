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
		return protocol.TypeBool
	case TypeByte:
		return protocol.TypeByte
	case TypeI16:
		return protocol.TypeI16
	case TypeI32:
		return protocol.TypeI32
	case TypeI64:
		return protocol.TypeI64
	case TypeDouble:
		return protocol.TypeDouble
	case TypeBinary:
		return protocol.TypeString
	case TypeList:
		return protocol.TypeList
	case TypeSet:
		return protocol.TypeSet
	case TypeMap:
		return protocol.TypeMap
	case TypeStruct:
		return protocol.TypeStruct
	}
	return protocol.TypeStop
}
