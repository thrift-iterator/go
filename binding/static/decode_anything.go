package static

import (
	"github.com/v2pro/wombat/generic"
	"reflect"
)

var byteArrayType = reflect.TypeOf(([]byte)(nil))

func dispatch(dstType reflect.Type, srcType reflect.Type) string {
	dstType = dstType.Elem()
	if dstType == byteArrayType {
		return "DecodeBinary"
	}
	if isEnumType(dstType) {
		return "DecodeEnum"
	}
	switch dstType.Kind() {
	case reflect.Slice:
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

func isEnumType(valType reflect.Type) bool {
	if valType.Kind() != reflect.Int64 {
		return false
	}
	_, hasStringMethod := valType.MethodByName("String")
	return hasStringMethod
}

var decodeAnything = generic.DefineFunc("DecodeAnything(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"dispatch", dispatch).
	Source(`
{{ $tmpl := dispatch .DT .ST }}
{{ $decode := expand $tmpl "DT" .DT "ST" .ST }}
{{$decode}}(dst, src)
`)
