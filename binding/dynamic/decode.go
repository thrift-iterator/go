package dynamic

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"unsafe"
)

var byteSliceType = reflect.TypeOf(([]byte)(nil))

func DecoderOf(valType reflect.Type) spi.ValDecoder {
	return &valDecoderAdapter{decoderOf("", valType)}
}

func decoderOf(prefix string, valType reflect.Type) internalDecoder {
	valType = valType.Elem()
	if byteSliceType == valType {
		return &binaryDecoder{}
	}
	switch valType.Kind() {
	case reflect.Bool:
		return &boolDecoder{}
	case reflect.Float64:
		return &float64Decoder{}
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
