package general

import "github.com/thrift-iterator/go/spi"

type generalMapDecoder struct {
}

func (decoder *generalMapDecoder) Decode(val interface{}, iter spi.Iterator) {
	*val.(*map[interface{}]interface{}) = readMap(iter).(map[interface{}]interface{})
}

func readMap(iter spi.Iterator) interface{} {
	keyType, elemType, length := iter.ReadMapHeader()
	keyReader := generalReaderOf(keyType)
	elemReader := generalReaderOf(elemType)
	generalMap := map[interface{}]interface{}{}
	for i := 0; i < length; i++ {
		key := keyReader(iter)
		elem := elemReader(iter)
		generalMap[key] = elem
	}
	return generalMap
}