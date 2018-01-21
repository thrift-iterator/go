package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
)

func dispatchEncode(dstType reflect.Type, srcType reflect.Type) string {
	if srcType == byteArrayType {
		return "EncodeBinary"
	}
	if isEnumType(srcType) {
		return "EncodeEnum"
	}
	switch srcType.Kind() {
	case reflect.Slice:
		return "EncodeSlice"
	case reflect.Map:
		return "EncodeMap"
	case reflect.Struct:
		return "EncodeStruct"
	case reflect.Ptr:
		return "EncodePointer"
	}
	return "EncodeSimpleValue"
}

var encodeAnything = generic.DefineFunc("EncodeAnything(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"dispatchEncode", dispatchEncode).
	Source(`
{{ $tmpl := dispatchEncode .DT .ST }}
{{ $encode := expand $tmpl "DT" .DT "ST" .ST }}
{{$encode}}(dst, src)
`)