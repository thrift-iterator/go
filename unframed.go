package thrifter

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"io"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/general"
)

type unframedDecoder struct {
	decodeFromReader bool
	cfg              *frozenConfig
	iter             spi.Iterator
}

type unframedEncoder struct {
	cfg    *frozenConfig
	stream spi.Stream
}

func (decoder *unframedDecoder) Decode(val interface{}) error {
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

func (decoder *unframedDecoder) DecodeMessage() (general.Message, error) {
	var msg general.Message
	err := decoder.Decode(&msg)
	return msg, err
}

func (decoder *unframedDecoder) DecodeMessageHeader() (protocol.MessageHeader, error) {
	var msgHeader protocol.MessageHeader
	err := decoder.Decode(&msgHeader)
	return msgHeader, err
}

func (decoder *unframedDecoder) DecodeMessageArguments() (general.Struct, error) {
	var msgArgs general.Struct
	err := decoder.Decode(&msgArgs)
	return msgArgs, err
}

func (decoder *unframedDecoder) Reset(reader io.Reader, buf []byte) {
	decoder.iter.Reset(reader, buf)
}

func (encoder *unframedEncoder) Encode(val interface{}) error {
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

func (encoder *unframedEncoder) EncodeMessage(msg general.Message) error {
	return encoder.Encode(msg)
}

func (encoder *unframedEncoder) EncodeMessageHeader(msgHeader protocol.MessageHeader) error {
	return encoder.Encode(msgHeader)
}

func (encoder *unframedEncoder) EncodeMessageArguments(msgArgs general.Struct) error {
	return encoder.Encode(msgArgs)
}

func (encoder *unframedEncoder) Reset(writer io.Writer) {
	encoder.stream.Reset(writer)
}

func (encoder *unframedEncoder) Buffer() []byte {
	return encoder.stream.Buffer()
}
