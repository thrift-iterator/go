package dynamic

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"unsafe"
	"github.com/thrift-iterator/go/protocol"
)

func EncoderOf(valType reflect.Type) spi.ValEncoder {
	return &valEncoderAdapter{encoderOf("", valType)}
}

func encoderOf(prefix string, valType reflect.Type) internalEncoder {
	if byteSliceType == valType {
		return &binaryEncoder{}
	}
	if isEnumType(valType) {
		return &int32Encoder{}
	}
	switch valType.Kind() {
	case reflect.String:
		return &stringEncoder{}
	case reflect.Bool:
		return &boolEncoder{}
	case reflect.Int8:
		return &int8Encoder{}
	case reflect.Uint8:
		return &uint8Encoder{}
	case reflect.Int16:
		return &int16Encoder{}
	case reflect.Uint16:
		return &uint16Encoder{}
	case reflect.Int32:
		return &int32Encoder{}
	case reflect.Uint32:
		return &uint32Encoder{}
	case reflect.Int64:
		return &int64Encoder{}
	case reflect.Uint64:
		return &uint64Encoder{}
	case reflect.Int:
		return &intEncoder{}
	case reflect.Uint:
		return &uintEncoder{}
	case reflect.Float32:
		return &float32Encoder{}
	case reflect.Float64:
		return &float64Encoder{}
	case reflect.Slice:
		return &sliceEncoder{
			sliceType: valType,
			elemType: valType.Elem(),
			elemEncoder: encoderOf(prefix + " [sliceElem]", valType.Elem()),
		}
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

func (encoder *unknownEncoder) thriftType() protocol.TType {
	return protocol.TypeStop
}
