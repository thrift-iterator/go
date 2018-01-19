package thrifter

import (
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/spi"
)

type msgDecoder struct {
}

func (decoder *msgDecoder) Decode(val interface{}, iter spi.Iterator) {
	msg, _ := val.(*protocol.Message)
	if msg == nil {
		iter.ReportError("MsgDecoder", "can only unmarshal protocol.Message")
		return
	}
	msgRead := iter.ReadMessage()
	if iter.Error() != nil {
		return
	}
	msg.Set(&msgRead)
}

var msgDecoderInstance = &msgDecoder{}