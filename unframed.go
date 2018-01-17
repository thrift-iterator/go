package thrifter

import (
	"github.com/thrift-iterator/go/protocol"
	"errors"
)

type unframedDecoder struct {
	iter Iterator
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