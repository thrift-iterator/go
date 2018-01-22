package static

import (
	"github.com/v2pro/wombat/generic"
	"reflect"
)

func dispatchDecode(dstType reflect.Type) string {
	dstType = dstType.Elem()
	if dstType == byteArrayType {
		return "DecodeBinary"
	}
	if isEnumType(dstType) {
		return "DecodeEnum"
	}
	switch dstType.Kind() {
	case reflect.Slice:
		if dstType.Elem().Kind() == reflect.Interface {
			return "DecodeSliceOfObject"
		}
		return "DecodeSlice"
	case reflect.Map:
		return "DecodeMap"
	case reflect.Struct:
		return "DecodeStruct"
	case reflect.Ptr:
		return "DecodePointer"
	}
	return "DecodeSimpleValue"
}

var decodeAnything = generic.DefineFunc("DecodeAnything(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"dispatchDecode", dispatchDecode).
	Source(`
{{ $tmpl := dispatchDecode .DT }}
{{ $decode := expand $tmpl "DT" .DT "ST" .ST }}
{{$decode}}(dst, src)
`)
