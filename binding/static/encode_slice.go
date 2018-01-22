package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
)

func init() {
	encodeAnything.ImportFunc(encodeSlice)
	encodeAnything.ImportFunc(encodeSliceOfObject)
}

var encodeSlice = generic.DefineFunc(
	"EncodeSlice(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(encodeAnything).
	Generators(
	"thriftType", func(srcType reflect.Type) int {
		_, ttype := dispatchEncode(srcType)
		return int(ttype)
	}).
	Source(`
{{ $encodeElem := expand "EncodeAnything" "DT" .DT "ST" (.ST|elem) }}
dst.WriteListHeader({{.ST|elem|thriftType}}, len(src))
for _, elem := range src {
	{{$encodeElem}}(dst, elem)
}
`)

var encodeSliceOfObject = generic.DefineFunc(
	"EncodeSliceOfObject(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Source(`
dst.WriteList(src)
`)