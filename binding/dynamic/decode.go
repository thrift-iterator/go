package dynamic

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
)

func DecoderOf(valType reflect.Type) spi.ValDecoder {
	return &binaryDecoder{}
}