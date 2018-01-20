package static

import (
	"github.com/v2pro/wombat/generic"
)

func init() {
	decodeAnything.ImportFunc(decodingPointer)
}

var decodingPointer = generic.DefineFunc(
	"DecodePointer(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(decodeAnything).
	Source(`
{{ $decode := expand "DecodeAnything" "DT" (.DT|elem) "ST" .ST }}
defDst := new({{ .DT|elem|elem|name }})
{{$decode}}(defDst, src)
*dst = defDst
return`)