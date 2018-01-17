package thrifter

import (
	"github.com/thrift-iterator/go/protocol"
	"errors"
)

type unframedDecoder struct {
	iter Iterator
}

type unframedEncoder struct {
	stream Stream
}

func (decoder *unframedDecoder) Decode(obj interface{}) error {
	msg, _ := obj.(*protocol.Message)
	if msg == nil {
		return errors.New("can only unmarshal protocol.Message")
	}
	msgRead := decoder.iter.ReadMessage()
	if decoder.iter.Error() != nil {
		return decoder.iter.Error()
	}
	msg.Set(&msgRead)
	return nil
}

func (encoder *unframedEncoder) Encode(obj interface{}) error {
	msg, isMsg := obj.(protocol.Message)
	if !isMsg {
		return errors.New("can only marshal protocol.Message")
	}
	encoder.stream.WriteMessage(msg)
	encoder.stream.Flush()
	if encoder.stream.Error() != nil {
		return encoder.stream.Error()
	}
	return nil
}