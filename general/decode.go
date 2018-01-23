package general

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

func ExportDecoders() map[reflect.Type]spi.ValDecoder {
	return map[reflect.Type]spi.ValDecoder {
		reflect.TypeOf((*[]interface{})(nil)): &generalListDecoder{},
	}
}

func generalReaderOf(ttype protocol.TType) func(iter spi.Iterator) interface{} {
	switch ttype {
	case protocol.TypeI64:
		return readInt64
	default:
		panic("unsupported type")
	}
}

func readInt64(iter spi.Iterator) interface{} {
	return iter.ReadInt64()
}