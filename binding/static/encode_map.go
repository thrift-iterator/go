package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
)

func init() {
	encodeAnything.ImportFunc(encodeMap)
}

var encodeMap = generic.DefineFunc(
	"EncodeMap(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(encodeAnything).
	Generators(
	"thriftType", func(srcType reflect.Type) int {
		_, ttype := dispatchEncode(srcType)
		return int(ttype)
	}).
	Source(`
{{ $encodeKey := expand "EncodeAnything" "DT" .DT "ST" (.ST|key) }}
{{ $encodeElem := expand "EncodeAnything" "DT" .DT "ST" (.ST|elem) }}
dst.WriteMapHeader({{.ST|key|thriftType}}, {{.ST|elem|thriftType}}, len(src))
for key, elem := range src {
	{{$encodeKey}}(dst, key)
	{{$encodeElem}}(dst, elem)
}`)