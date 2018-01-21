package dynamic

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"unsafe"
)

func EncoderOf(valType reflect.Type) spi.ValEncoder {
	return &valEncoderAdapter{encoderOf("", valType)}
}

func encoderOf(prefix string, valType reflect.Type) internalEncoder {
	if isEnumType(valType) {
		return &int32Encoder{}
	}
	switch valType.Kind() {
	case reflect.Bool:
		return &boolEncoder{}
	case reflect.Float64:
		return &float64Encoder{}
	}
	return &unknownEncoder{prefix, valType}
}

type unknownEncoder struct {
	prefix  string
	valType reflect.Type
}

func (encoder *unknownEncoder) encode(ptr unsafe.Pointer, stream spi.Stream) {
	stream.ReportError("decode "+encoder.prefix, "do not know how to encode "+encoder.valType.String())
}
