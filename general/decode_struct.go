package general

import (
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

type generalStructDecoder struct {
}

func (decoder *generalStructDecoder) Decode(val interface{}, iter spi.Iterator) {
	obj := *val.(*map[protocol.FieldId]interface{})
	if obj == nil {
		obj = map[protocol.FieldId]interface{}{}
		*val.(*map[protocol.FieldId]interface{}) = obj
	}
	iter.ReadStructHeader()
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.TypeStop {
			return
		}
		generalReader := generalReaderOf(fieldType)
		obj[fieldId] = generalReader(iter)
	}
}