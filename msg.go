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

type msgHeaderDecoder struct {
}

func (decoder *msgHeaderDecoder) Decode(val interface{}, iter spi.Iterator) {
	msgHeader := val.(*protocol.MessageHeader)
	msgHeaderRead := iter.ReadMessageHeader()
	if iter.Error() != nil {
		return
	}
	msgHeader.Set(&msgHeaderRead)
}

var msgHeaderDecoderInstance = &msgHeaderDecoder{}

type msgEncoder struct {
}

func (encoder *msgEncoder) Encode(val interface{}, stream spi.Stream) {
	msg := val.(protocol.Message)
	stream.WriteMessage(msg)
}

var msgEncoderInstance = &msgEncoder{}

type msgHeaderEncoder struct {
}

func (encoder *msgHeaderEncoder) Encode(val interface{}, stream spi.Stream) {
	msgHeader := val.(protocol.MessageHeader)
	stream.WriteMessageHeader(msgHeader)
}

var msgHeaderEncoderInstance = &msgHeaderEncoder{}