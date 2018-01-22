package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go/protocol"
)

func dispatchEncode(srcType reflect.Type) (string, protocol.TType) {
	if srcType == byteArrayType {
		return "EncodeBinary", protocol.TypeString
	}
	if isEnumType(srcType) {
		return "EncodeEnum", protocol.TypeI32
	}
	switch srcType.Kind() {
	case reflect.Slice:
		if srcType.Elem().Kind() == reflect.Interface {
			return "EncodeSliceOfObject", protocol.TypeList
		}
		return "EncodeSlice", protocol.TypeList
	case reflect.Map:
		return "EncodeMap", protocol.TypeMap
	case reflect.Struct:
		return "EncodeStruct", protocol.TypeStruct
	case reflect.Ptr:
		_, ttype := dispatchEncode(srcType.Elem())
		return "EncodePointer", ttype
	}
	return "EncodeSimpleValue", thriftTypeMap[srcType.Kind()]
}

var encodeAnything = generic.DefineFunc("EncodeAnything(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"dispatchEncode", func(srcType reflect.Type) string {
		encode, _ := dispatchEncode(srcType)
		return encode
	}).
	Source(`
{{ $tmpl := dispatchEncode .ST }}
{{ $encode := expand $tmpl "DT" .DT "ST" .ST }}
{{$encode}}(dst, src)
`)