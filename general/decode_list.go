package general

import "github.com/thrift-iterator/go/spi"

type generalListDecoder struct {
}

func (decoder *generalListDecoder) Decode(val interface{}, iter spi.Iterator) {
	obj := val.(*[]interface{})
	elemType, length := iter.ReadListHeader()
	generalReader := generalReaderOf(elemType)
	for i := 0; i < length; i++ {
		*obj = append(*obj, generalReader(iter))
	}
}