package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
)

func init() {
	decodeAnything.ImportFunc(decodeSlice)
	decodeAnything.ImportFunc(decodeSliceOfObject)
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
_, length := src.ReadListHeader()
for i := 0; i < length; i++ {
	elem := new({{.DT|elem|elem|name}})
	{{$decodeElem}}(elem, src)
	*dst = append(*dst, *elem)
}`)

var decodeSliceOfObject = generic.DefineFunc(
	"DecodeSliceOfObject(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Source(`
elemType, length := src.ReadListHeader()
for i := 0; i < length; i++ {
	*dst = append(*dst, src.Read(elemType))
}`)