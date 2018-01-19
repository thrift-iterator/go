package dynamic

import "github.com/thrift-iterator/go/spi"

type binaryDecoder struct {
}

func (decoder *binaryDecoder) Decode(val interface{}, iter spi.Iterator) {
	*val.(*[]byte) = iter.ReadBinary()
}