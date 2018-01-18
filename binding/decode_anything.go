package binding

import (
	"github.com/v2pro/wombat/generic"
	"reflect"
)

var byteArrayType = reflect.TypeOf(([]byte)(nil))

func dispatch(dstType reflect.Type, srcType reflect.Type) string {
	if dstType.Kind() != reflect.Map {
		dstType = dstType.Elem()
	}
	if dstType == byteArrayType {
		return "DecodeBinary"
	}
	switch dstType.Kind() {
	case reflect.Slice:
		return "DecodeSlice"
	case reflect.Map:
		return "DecodeMap"
	case reflect.Struct:
		return "DecodeStruct"
	}
	return "DecodeSimpleValue"
}

var decodeAnything = generic.DefineFunc("DecodeAnything(dst interface{}, src interface{})").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"dispatch", dispatch).
	Source(`
{{ $tmpl := dispatch .DT .ST }}
{{ $decode := expand $tmpl "DT" .DT "ST" .ST }}
{{$decode}}(dst.({{.DT|name}}), src.({{.ST|name}}))
`)
