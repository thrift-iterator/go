package thrifter

import (
	"github.com/thrift-iterator/go/protocol"
	"errors"
	"reflect"
)

type unframedDecoder struct {
	cfg *frozenConfig
	iter Iterator
}

type unframedEncoder struct {
	stream Stream
}

func (decoder *unframedDecoder) Decode(obj interface{}) error {
	valDecoder := decoder.cfg.decoders[reflect.TypeOf(obj)]
	if valDecoder == nil {
		valDecoder = msgDecoderInstance
	}
	valDecoder.Decode(obj, decoder.iter)
	if decoder.iter.Error() != nil {
		return decoder.iter.Error()
	}
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