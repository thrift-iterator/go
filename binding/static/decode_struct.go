package static

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go/protocol"
	"strconv"
	"strings"
)

func init() {
	decodeAnything.ImportFunc(decodeStruct)
}

var decodeStruct = generic.DefineFunc(
	"DecodeStruct(dst DT, src ST)").
	Param("DT", "the dst type to copy into").
	Param("ST", "the src type to copy from").
	ImportFunc(decodeAnything).
	Generators(
	"calcBindings", func(dstType, srcType reflect.Type) interface{} {
		bindings := []interface{}{}
		for i := 0; i < dstType.NumField(); i++ {
			dstField := dstType.Field(i)
			srcFieldId := protocol.FieldId(0)
			thriftTag := dstField.Tag.Get("thrift")
			if thriftTag != "" {
				parts := strings.Split(thriftTag, ",")
				if len(parts) >= 2 {
					fieldId, err := strconv.Atoi(parts[1])
					if err != nil {
						panic("thrift tag must be integer")
					}
					srcFieldId = protocol.FieldId(fieldId)
				}
			}
			bindings = append(bindings, map[string]interface{}{
				"srcFieldId": srcFieldId,
				"srcType": srcType,
				"dstFieldName": dstField.Name,
				"dstFieldType": reflect.PtrTo(dstField.Type),
			})
		}
		return bindings
	},
	"assignDecode", func(binding map[string]interface{}, decodeFuncName string) string {
		binding["decode"] = decodeFuncName
		return ""
	}).
	Source(`
{{ $bindings := calcBindings (.DT|elem) .ST }}
{{ range $_, $binding := $bindings}}
	{{ $decode := expand "DecodeAnything" "DT" $binding.dstFieldType "ST" $binding.srcType }}
	{{ assignDecode $binding $decode }}
{{ end }}
src.ReadStructHeader()
for {
	fieldType, fieldId := src.ReadStructField()
	if fieldType == 0 {
		return
	}
	switch fieldId {
		{{ range $_, $binding := $bindings }}
			case {{ $binding.srcFieldId }}:
				{{$binding.decode}}(&dst.{{$binding.dstFieldName}}, src)
		{{ end }}
		default:
			src.Discard(fieldType)
	}
}`)