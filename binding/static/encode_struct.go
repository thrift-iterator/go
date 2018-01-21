package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
)

func init() {
	encodeAnything.ImportFunc(encodeStruct)
}

var encodeStruct = generic.DefineFunc(
	"EncodeStruct(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(encodeAnything).
	Generators(
	"calcBindings", calcBindings,
	"assignEncode", func(binding map[string]interface{}, encodeFuncName string) string {
		binding["encode"] = encodeFuncName
		return ""
	},
	"thriftType", func(srcType reflect.Type) int {
		_, ttype := dispatchEncode(srcType)
		return int(ttype)
	}).
	Source(`
{{ $bindings := calcBindings .ST }}
dst.WriteStructHeader()
{{ range $_, $binding := $bindings}}
	{{ $encode := expand "EncodeAnything" "DT" $.DT "ST" $binding.fieldType }}
	dst.WriteStructField({{$binding.fieldType|thriftType}}, {{$binding.fieldId}})
	{{$encode}}(dst, &src.{{$binding.fieldName}})
{{ end }}
dst.WriteStructFieldStop()
`)
