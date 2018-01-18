package binding

import (
	"github.com/v2pro/wombat/generic"
	"reflect"
)

var decodeSimpleValue = generic.DefineFunc(
	"DecodeSimpleValue(err *error, dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	Generators(
	"opFuncName", func(typ reflect.Type) string {
		return simpleValueMap[typ.Kind()]
	}).
	Source(`
*dst = src.Read{{.DT|elem|opFuncName}}()
	`)