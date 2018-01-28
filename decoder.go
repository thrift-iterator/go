package thrifter

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"io"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/general"
)

type Decoder struct {
	decodeFromReader bool
	cfg              *frozenConfig
	iter             spi.Iterator
}

func (decoder *Decoder) Decode(val interface{}) error {
	cfg := decoder.cfg
	valType := reflect.TypeOf(val)
	valDecoder := cfg.getGenDecoder(valType)
	if valDecoder == nil {
		valDecoder = cfg.decoderOf(decoder.decodeFromReader, valType)
		cfg.addGenDecoder(valType, valDecoder)
	}
	valDecoder.Decode(val, decoder.iter)
	if decoder.iter.Error() != nil {
		return decoder.iter.Error()
	}
	return nil
}

func (decoder *Decoder) DecodeMessage() (general.Message, error) {
	var msg general.Message
	err := decoder.Decode(&msg)
	return msg, err
}

func (decoder *Decoder) DecodeMessageHeader() (protocol.MessageHeader, error) {
	var msgHeader protocol.MessageHeader
	err := decoder.Decode(&msgHeader)
	return msgHeader, err
}

func (decoder *Decoder) DecodeMessageArguments() (general.Struct, error) {
	var msgArgs general.Struct
	err := decoder.Decode(&msgArgs)
	return msgArgs, err
}

func (decoder *Decoder) Reset(reader io.Reader, buf []byte) {
	decoder.iter.Reset(reader, buf)
}

func (encoder *Encoder) Encode(val interface{}) error {
	cfg := encoder.cfg
	valType := reflect.TypeOf(val)
	valEncoder := cfg.getGenEncoder(valType)
	if valEncoder == nil {
		valEncoder = cfg.encoderOf(valType)
		cfg.addGenEncoder(valType, valEncoder)
	}
	valEncoder.Encode(val, encoder.stream)
	encoder.stream.Flush()
	if encoder.stream.Error() != nil {
		return encoder.stream.Error()
	}
	return nil
}
