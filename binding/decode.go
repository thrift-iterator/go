package binding

import "github.com/v2pro/wombat/generic"

var Decode = generic.DefineFunc("Decode(dst interface{}, src interface{})").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(decodeAnything).
	Source(`
{{ $decode := expand "DecodeAnything" "DT" .DT "ST" .ST }}
{{$decode}}(dst.({{.DT|name}}), src.({{.ST|name}}))
`)