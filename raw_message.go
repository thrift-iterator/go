package thrifter

import (
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

// RawMessage works like json.RawMessage, keep original buffer as it is
type RawMessage struct {
	ThriftType protocol.TType
	Buffer     []byte
}

type rawStructDecoder struct {
}

func (decoder *rawStructDecoder) Decode(val interface{}, iter spi.Iterator) {
	rawStruct := *val.(*map[protocol.FieldId]RawMessage)
	iter.ReadStructHeader()
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.TypeStop {
			return
		}
		skipped := iter.Skip(fieldType, nil)
		// make a copy, so we can safely reference the buffer later
		rawMessage := RawMessage{
			Buffer:     append(([]byte)(nil), skipped...),
			ThriftType: fieldType,
		}
		rawStruct[fieldId] = rawMessage
	}
}

var rawStructDecoderInstance = &rawStructDecoder{}

type rawStructEncoder struct {
}

func (encoder *rawStructEncoder) Encode(val interface{}, stream spi.Stream) {
	rawStruct := val.(map[protocol.FieldId]RawMessage)
	stream.WriteStructHeader()
	for key, elem := range rawStruct {
		stream.WriteStructField(elem.ThriftType, key)
		stream.Write(elem.Buffer)
	}
	stream.WriteStructFieldStop()
}

func (encoder *rawStructEncoder) ThriftType() protocol.TType {
	return protocol.TypeStruct
}

var rawStructEncoderInstance = &rawStructEncoder{}
