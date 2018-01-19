package static

import (
	"github.com/v2pro/wombat/generic"
	"reflect"
)

func init() {
	decodeAnything.ImportFunc(decodeSimpleValue)
}

var decodeSimpleValue = generic.DefineFunc(
	"DecodeSimpleValue(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"opFuncName", func(typ reflect.Type) string {
		return simpleValueMap[typ.Kind()]
	}).
	Source(`
*dst = src.Read{{.DT|elem|opFuncName}}()
	`)