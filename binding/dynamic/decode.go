package dynamic

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"unsafe"
)

var byteSliceType = reflect.TypeOf(([]byte)(nil))

func DecoderOf(valType reflect.Type) spi.ValDecoder {
	return &valDecoderAdapter{decoderOf("", valType.Elem())}
}

func decoderOf(prefix string, valType reflect.Type) internalDecoder {
	if byteSliceType == valType {
		return &binaryDecoder{}
	}
	switch valType.Kind() {
	case reflect.Bool:
		return &boolDecoder{}
	case reflect.Float64:
		return &float64Decoder{}
	case reflect.Int8:
		return &int8Decoder{}
	case reflect.Uint8:
		return &uint8Decoder{}
	case reflect.Int16:
		return &int16Decoder{}
	case reflect.Uint16:
		return &uint16Decoder{}
	case reflect.Int32:
		return &int32Decoder{}
	case reflect.Uint32:
		return &uint32Decoder{}
	case reflect.Int64:
		return &int64Decoder{}
	case reflect.Uint64:
		return &uint64Decoder{}
	case reflect.String:
		return &stringDecoder{}
	case reflect.Slice:
		return &sliceDecoder{
			elemType: valType.Elem(),
			sliceType: valType,
			elemDecoder: decoderOf(prefix + " [sliceElem]", valType.Elem()),
		}
	}
	return &unknownDecoder{prefix, valType}
}

type unknownDecoder struct {
	prefix  string
	valType reflect.Type
}

func (decoder *unknownDecoder) decode(ptr unsafe.Pointer, iterator spi.Iterator) {
	iterator.ReportError("decode " + decoder.prefix, "do not know how to decode "+decoder.valType.String())
}
