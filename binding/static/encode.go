package static

import "github.com/v2pro/wombat/generic"

var Encode = generic.DefineFunc("Encode(dst interface{}, src interface{})").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(encodeAnything).
	Source(`
{{ $decode := expand "EncodeAnything" "DT" .DT "ST" .ST }}
{{$decode}}(dst.({{.DT|name}}), src.({{.ST|name}}))
`)