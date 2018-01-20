package thrifter

import (
	"github.com/thrift-iterator/go/protocol"
	"errors"
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"io"
)

type unframedDecoder struct {
	cfg  *frozenConfig
	iter spi.Iterator
}

type unframedEncoder struct {
	stream spi.Stream
}

func (decoder *unframedDecoder) Decode(val interface{}) error {
	cfg := decoder.cfg
	valType := reflect.TypeOf(val)
	valDecoder := cfg.getDecoderFromCache(valType)
	if valDecoder == nil {
		valDecoder = cfg.decoderOf(true, valType)
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
