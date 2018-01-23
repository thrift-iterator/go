package general

import (
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

type messageDecoder struct {
}

func (decoder *messageDecoder) Decode(val interface{}, iter spi.Iterator) {
	*val.(*protocol.Message) = protocol.Message{
		MessageHeader: iter.ReadMessageHeader(),
		Arguments:     readStruct(iter).(map[protocol.FieldId]interface{}),
	}
}

type messageHeaderDecoder struct {
}

func (decoder *messageHeaderDecoder) Decode(val interface{}, iter spi.Iterator) {
	*val.(*protocol.MessageHeader) = iter.ReadMessageHeader()
}