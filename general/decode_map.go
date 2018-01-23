package general

import "github.com/thrift-iterator/go/spi"

type generalMapDecoder struct {
}

func (decoder *generalMapDecoder) Decode(val interface{}, iter spi.Iterator) {
	keyType, elemType, length := iter.ReadMapHeader()
	keyReader := generalReaderOf(keyType)
	elemReader := generalReaderOf(elemType)
	obj := *val.(*map[interface{}]interface{})
	if obj == nil {
		obj = map[interface{}]interface{}{}
		*val.(*map[interface{}]interface{}) = obj
	}
	for i := 0; i < length; i++ {
		key := keyReader(iter)
		elem := elemReader(iter)
		obj[key] = elem
	}
}