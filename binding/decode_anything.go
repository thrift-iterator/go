package binding

import (
	"github.com/v2pro/wombat/generic"
	"reflect"
)

var pByteArrayType = reflect.TypeOf((*[]byte)(nil))

func dispatch(dstType reflect.Type, srcType reflect.Type) string {
	if dstType == pByteArrayType {
		return "DecodeBinary"
	}
	return "DecodeSimpleValue"
}

var DecodeAnything = generic.DefineFunc("DecodeAnything(dst interface{}, src interface{})").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"dispatch", dispatch).
	Source(`
{{ $tmpl := dispatch .DT .ST }}
{{ $decode := expand $tmpl "DT" .DT "ST" .ST }}
{{$decode}}(dst.({{.DT|name}}), src.({{.ST|name}}))
`)