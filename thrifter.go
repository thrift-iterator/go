package thrifter

import (
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol"
)

type Protocol int

var ProtocolBinary Protocol = 1

type Iterator interface {
	ReadStructCB(func(fieldType protocol.TType, fieldId protocol.FieldId))
	ReadBool() bool
	ReadInt8() int8
	ReadUInt8() uint8
	ReadInt16() int16
	ReadUInt16() uint16
	ReadInt32() int32
	ReadUInt32() uint32
	ReadInt64() int64
	ReadUInt64() uint64
	ReadFloat64() float64
	ReadString() string
	ReadBinary() []byte
}

type Config struct {
	Protocol Protocol
}

type API interface {
	NewIterator(buf []byte) Iterator
}

type frozenConfig struct {
	protocol Protocol
}

func (cfg Config) Froze() API {
	api := &frozenConfig{protocol: cfg.Protocol}
	return api
}

func (cfg *frozenConfig) NewIterator(buf []byte) Iterator {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewIterator(buf)
	}
	panic("unsupported protocol")
}

var DefaultConfig = Config{Protocol: ProtocolBinary}.Froze()

func NewIterator(buf []byte) Iterator {
	return DefaultConfig.NewIterator(buf)
}
