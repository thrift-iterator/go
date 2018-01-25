package raw

import "github.com/thrift-iterator/go/protocol"

type Buffer []byte

type StructField struct {
	Buffer Buffer
	Type protocol.TType
}

type Struct map[protocol.FieldId]StructField

type List struct {
	ElementType protocol.TType
	Elements []Buffer
}

type Map struct {
	KeyType protocol.TType
	ElementType protocol.TType
	Entries map[interface{}]Buffer
}