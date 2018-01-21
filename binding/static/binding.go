package static

import (
	"reflect"
	"github.com/thrift-iterator/go/protocol"
)

var byteArrayType = reflect.TypeOf(([]byte)(nil))

var simpleValueMap = map[reflect.Kind]string{
	reflect.Int: "Int",
	reflect.Int8: "Int8",
	reflect.Int16: "Int16",
	reflect.Int32: "Int32",
	reflect.Int64: "Int64",
	reflect.Uint: "Uint",
	reflect.Uint8: "Uint8",
	reflect.Uint16: "Uint16",
	reflect.Uint32: "Uint32",
	reflect.Uint64: "Uint64",
	reflect.Float32: "Float32",
	reflect.Float64: "Float64",
	reflect.String: "String",
	reflect.Bool: "Bool",
}

var thriftTypeMap = map[reflect.Kind]protocol.TType {
	reflect.Int: protocol.TypeI64,
	reflect.Int8: protocol.TypeI08,
	reflect.Int16: protocol.TypeI16,
	reflect.Int32: protocol.TypeI32,
	reflect.Int64: protocol.TypeI64,
	reflect.Uint: protocol.TypeI64,
	reflect.Uint8: protocol.TypeI08,
	reflect.Uint16: protocol.TypeI16,
	reflect.Uint32: protocol.TypeI32,
	reflect.Uint64: protocol.TypeI64,
	reflect.Float32: protocol.TypeDouble,
	reflect.Float64: protocol.TypeDouble,
	reflect.String: protocol.TypeString,
	reflect.Bool: protocol.TypeBool,
}

func isEnumType(valType reflect.Type) bool {
	if valType.Kind() != reflect.Int64 {
		return false
	}
	_, hasStringMethod := valType.MethodByName("String")
	return hasStringMethod
}