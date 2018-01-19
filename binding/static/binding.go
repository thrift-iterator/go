package static

import "reflect"

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