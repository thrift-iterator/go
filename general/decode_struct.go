package general

import (
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

type generalStructDecoder struct {
}

func (decoder *generalStructDecoder) Decode(val interface{}, iter spi.Iterator) {
	*val.(*map[protocol.FieldId]interface{}) = readStruct(iter).(map[protocol.FieldId]interface{})
}

func readStruct(iter spi.Iterator) interface{} {
	generalStruct := map[protocol.FieldId]interface{}{}
	iter.ReadStructHeader()
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.TypeStop {
			return generalStruct
		}
		generalReader := generalReaderOf(fieldType)
		generalStruct[fieldId] = generalReader(iter)
	}
}