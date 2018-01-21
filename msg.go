package thrifter

import (
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/spi"
)

type msgDecoder struct {
}

func (decoder *msgDecoder) Decode(val interface{}, iter spi.Iterator) {
	msg := val.(*protocol.Message)
	msgRead := iter.ReadMessage()
	if iter.Error() != nil {
		return
	}
	msg.Set(&msgRead)
}

var msgDecoderInstance = &msgDecoder{}

type msgEncoder struct {
}

func (encoder *msgEncoder) Encode(val interface{}, stream spi.Stream) {
	msg := val.(protocol.Message)
	stream.WriteMessage(msg)
}

var msgEncoderInstance = &msgEncoder{}