package general

import "github.com/thrift-iterator/go/spi"

type generalMapDecoder struct {
}

func (decoder *generalMapDecoder) Decode(val interface{}, iter spi.Iterator) {
	*val.(*Map) = readMap(iter).(Map)
}

func readMap(iter spi.Iterator) interface{} {
	keyType, elemType, length := iter.ReadMapHeader()
	keyReader := ReaderOf(keyType)
	elemReader := ReaderOf(elemType)
	generalMap := Map{}
	for i := 0; i < length; i++ {
		key := keyReader(iter)
		elem := elemReader(iter)
		generalMap[key] = elem
	}
	return generalMap
}