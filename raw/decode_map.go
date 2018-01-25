package raw

import (
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/general"
)

type rawMapDecoder struct {
}

func (decoder *rawMapDecoder) Decode(val interface{}, iter spi.Iterator) {
	keyType, elemType, length := iter.ReadMapHeader()
	entries := make(map[interface{}][]byte, length)
	generalKeyReader := general.ReaderOf(keyType)
	for i := 0; i < length; i++ {
		key := generalKeyReader(iter)
		elemBuf := iter.Skip(elemType, nil)
		entries[key] = elemBuf
	}
	obj := val.(*Map)
	obj.KeyType = keyType
	obj.ElementType = elemType
	obj.Entries = entries
}