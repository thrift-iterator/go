package general

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

func ExportEncoders() map[reflect.Type]spi.ValEncoder{
	return map[reflect.Type]spi.ValEncoder {
		reflect.TypeOf(([]interface{})(nil)): &generalListEncoder{},
	}
}

func generalWriterOf(sample interface{}) (protocol.TType, func(val interface{}, stream spi.Stream)) {
	switch sample.(type) {
	case int64:
		return protocol.TypeI64, writeInt64
	default:
		panic("unsupported type")
	}
}

func writeInt64(val interface{}, stream spi.Stream) {
	stream.WriteInt64(val.(int64))
}