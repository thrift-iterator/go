package binding

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
)

func init() {
	decodeAnything.ImportFunc(decodeSlice)
}

var decodeSlice = generic.DefineFunc(
	"DecodeSlice(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(decodeAnything).
	Generators(
	"ptrSliceElem", func(typ reflect.Type) reflect.Type {
		return reflect.PtrTo(typ.Elem().Elem())
	}).
	Source(`
{{ $decodeElem := expand "DecodeAnything" "DT" (.DT|ptrSliceElem) "ST" .ST }}
originalLen := len(*dst)
_, length := src.ReadListHeader()
for i := 0; i < length; i++ {
	if i < originalLen {
		elem := &(*dst)[i]
		{{$decodeElem}}(elem, src)
	} else {
		elem := new({{.DT|elem|elem|name}})
		{{$decodeElem}}(elem, src)
		*dst = append(*dst, *elem)
	}
}`)