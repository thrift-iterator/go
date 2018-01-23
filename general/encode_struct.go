package general

import (
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

type generalStructEncoder struct {
}

func (encoder *generalStructEncoder) Encode(val interface{}, stream spi.Stream) {
	obj := val.(map[protocol.FieldId]interface{})
	stream.WriteStructHeader()
	for fieldId, elem := range obj {
		fieldType, generalWriter := generalWriterOf(elem)
		stream.WriteStructField(fieldType, fieldId)
		generalWriter(elem, stream)
	}
	stream.WriteStructFieldStop()
}
