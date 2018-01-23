package thrifter

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"io"
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
	valDecoder := cfg.GetDecoder(valType.String())
	if valDecoder == nil {
		valDecoder = cfg.decoderOf(decoder.decodeFromReader, valType)
		cfg.addDecoderToCache(valType, valDecoder)
	}
	valDecoder.Decode(val, decoder.iter)
	if decoder.iter.Error() != nil {
		return decoder.iter.Error()
	}
	return nil
}

func (decoder *unframedDecoder) Reset(reader io.Reader, buf []byte) {
	decoder.iter.Reset(reader, buf)
}

func (encoder *unframedEncoder) Encode(val interface{}) error {
	cfg := encoder.cfg
	valType := reflect.TypeOf(val)
	valEncoder := cfg.GetEncoder(valType.String())
	if valEncoder == nil {
		valEncoder = cfg.encoderOf(valType)
		cfg.addEncoderToCache(valType, valEncoder)
	}
	valEncoder.Encode(val, encoder.stream)
	encoder.stream.Flush()
	if encoder.stream.Error() != nil {
		return encoder.stream.Error()
	}
	return nil
}

func (encoder *unframedEncoder) Reset(writer io.Writer) {
	encoder.stream.Reset(writer)
}

func (encoder *unframedEncoder) Buffer() []byte {
	return encoder.stream.Buffer()
}
