package thrifter

import (
	"io"
	"encoding/binary"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/spi"
	"reflect"
)

type framedDecoder struct {
	cfg               *frozenConfig
	reader            io.Reader
	iter              spi.Iterator
	tmp               []byte
	shouldDecodeFrame bool
}

func (decoder *framedDecoder) Decode(val interface{}) error {
	if decoder.shouldDecodeFrame {
		if len(decoder.tmp) < 4 {
			decoder.tmp = make([]byte, 4)
		}
		tmp := decoder.tmp[:4]
		_, err := io.ReadFull(decoder.reader, tmp)
		if err != nil {
			return err
		}
		size := int(binary.BigEndian.Uint32(tmp))
		if len(decoder.tmp) < size {
			decoder.tmp = make([]byte, size)
		}
		tmp = decoder.tmp[:size]
		_, err = io.ReadFull(decoder.reader, tmp)
		if err != nil {
			return err
		}
		decoder.iter.Reset(nil, tmp)
		_, isMsg := val.(*protocol.Message)
		if !isMsg {
			decoder.shouldDecodeFrame = false
		}
	} else {
		decoder.shouldDecodeFrame = true
	}
	cfg := decoder.cfg
	valType := reflect.TypeOf(val)
	valDecoder := cfg.getGenDecoder(valType)
	if valDecoder == nil {
		valDecoder = cfg.decoderOf(true, valType)
		cfg.addGenDecoder(valType, valDecoder)
	}
	valDecoder.Decode(val, decoder.iter)
	return decoder.iter.Error()
}

func (decoder *framedDecoder) DecodeMessage() (protocol.Message, error) {
	var msg protocol.Message
	err := decoder.Decode(&msg)
	return msg, err
}

func (decoder *framedDecoder) DecodeMessageHeader() (protocol.MessageHeader, error) {
	var msgHeader protocol.MessageHeader
	err := decoder.Decode(&msgHeader)
	return msgHeader, err
}

func (decoder *framedDecoder) DecodeMessageArguments() (map[protocol.FieldId]interface{}, error) {
	var msgArgs map[protocol.FieldId]interface{}
	err := decoder.Decode(&msgArgs)
	return msgArgs, err
}

func (decoder *framedDecoder) Reset(reader io.Reader, buf []byte) {
	decoder.reader = reader
}

type framedEncoder struct {
	cfg               *frozenConfig
	writer            io.Writer
	stream            spi.Stream
	shouldEncodeFrame bool
}

func (encoder *framedEncoder) Encode(val interface{}) error {
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
	if _, isMsg := val.(*protocol.Message); isMsg {
		encoder.shouldEncodeFrame = true
	}
	if encoder.shouldEncodeFrame {
		buf := encoder.stream.Buffer()
		size := len(buf)
		_, err := encoder.writer.Write([]byte{
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		})
		if err != nil {
			return err
		}
		_, err = encoder.writer.Write(buf)
		if err != nil {
			return err
		}
		encoder.shouldEncodeFrame = false
	} else {
		encoder.shouldEncodeFrame = true
	}
	return nil
}

func (encoder *framedEncoder) EncodeMessage(msg protocol.Message) error {
	return encoder.Encode(msg)
}

func (encoder *framedEncoder) EncodeMessageHeader(msgHeader protocol.MessageHeader) error {
	return encoder.Encode(msgHeader)
}

func (encoder *framedEncoder) EncodeMessageArguments(msgArgs map[protocol.FieldId]interface{}) error {
	return encoder.Encode(msgArgs)
}

func (encoder *framedEncoder) Reset(writer io.Writer) {
	encoder.writer = writer
}

func (encoder *framedEncoder) Buffer() []byte {
	return encoder.stream.Buffer()
}
