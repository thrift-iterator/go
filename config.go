package thrifter

import (
	"unsafe"
	"reflect"
	"io"
	"github.com/thrift-iterator/go/protocol/sbinary"
	"github.com/thrift-iterator/go/protocol/compact"
	"sync/atomic"
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol"
	"errors"
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/binding/static"
	"github.com/thrift-iterator/go/binding/dynamic"
)

type frozenConfig struct {
	protocol       Protocol
	decoderCache   unsafe.Pointer
	encoderCache   unsafe.Pointer
	isFramed       bool
	dynamicCodegen bool
}

func (cfg Config) Froze() API {
	api := &frozenConfig{
		protocol: cfg.Protocol,
		isFramed: cfg.IsFramed,
		dynamicCodegen: cfg.DynamicCodegen,
	}
	atomic.StorePointer(&api.decoderCache, unsafe.Pointer(&map[reflect.Type]spi.ValDecoder{}))
	atomic.StorePointer(&api.encoderCache, unsafe.Pointer(&map[reflect.Type]spi.ValEncoder{}))
	return api
}

func (cfg *frozenConfig) addDecoderToCache(cacheKey reflect.Type, decoder spi.ValDecoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.decoderCache)
		cache := *(*map[reflect.Type]spi.ValDecoder)(ptr)
		copied := map[reflect.Type]spi.ValDecoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = decoder
		done = atomic.CompareAndSwapPointer(&cfg.decoderCache, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) addEncoderToCache(cacheKey reflect.Type, encoder spi.ValEncoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.encoderCache)
		cache := *(*map[reflect.Type]spi.ValEncoder)(ptr)
		copied := map[reflect.Type]spi.ValEncoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = encoder
		done = atomic.CompareAndSwapPointer(&cfg.encoderCache, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) getDecoderFromCache(cacheKey reflect.Type) spi.ValDecoder {
	ptr := atomic.LoadPointer(&cfg.decoderCache)
	cache := *(*map[reflect.Type]spi.ValDecoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) getEncoderFromCache(cacheKey reflect.Type) spi.ValEncoder {
	ptr := atomic.LoadPointer(&cfg.encoderCache)
	cache := *(*map[reflect.Type]spi.ValEncoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) NewStream(writer io.Writer, buf []byte) spi.Stream {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewStream(writer, buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) NewIterator(reader io.Reader, buf []byte) spi.Iterator {
	switch cfg.protocol {
	case ProtocolBinary:
		if reader != nil {
			return sbinary.NewIterator(reader, buf)
		}
		return binary.NewIterator(buf)
	case ProtocolCompact:
		return compact.NewIterator(buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) WillDecodeFromBuffer(samples ...interface{}) {
	for _, sample := range samples {
		cfg.staticDecoderOf(false, reflect.TypeOf(sample))
	}
}

func (cfg *frozenConfig) decoderOf(decodeFromReader bool, valType reflect.Type) spi.ValDecoder {
	if valType == reflect.TypeOf((*protocol.Message)(nil)) {
		return msgDecoderInstance
	}
	if cfg.dynamicCodegen {
		return dynamic.DecoderOf(valType)
	}
	return cfg.staticDecoderOf(decodeFromReader, valType)
}

func (cfg *frozenConfig) staticDecoderOf(decodeFromReader bool, valType reflect.Type) spi.ValDecoder {
	iteratorType := reflect.TypeOf((*binary.Iterator)(nil))
	if decodeFromReader {
		iteratorType = reflect.TypeOf((*sbinary.Iterator)(nil))
	}
	if cfg.protocol == ProtocolCompact {
		iteratorType = reflect.TypeOf((*compact.Iterator)(nil))
	}
	funcObj := generic.Expand(static.Decode,
		"ST", iteratorType,
		"DT", valType)
	f := funcObj.(func(interface{}, interface{}))
	return &funcDecoder{f}
}

type funcDecoder struct {
	f func(dst interface{}, src interface{})
}

func (decoder *funcDecoder) Decode(val interface{}, iter spi.Iterator) {
	decoder.f(val, iter)
}

func (cfg *frozenConfig) Unmarshal(buf []byte, val interface{}) error {
	valType := reflect.TypeOf(val)
	decoder := cfg.getDecoderFromCache(valType)
	if decoder == nil {
		decoder = cfg.decoderOf(false, valType)
		cfg.addDecoderToCache(valType, decoder)
	}
	if buf == nil {
		return errors.New("empty input")
	}
	if cfg.isFramed {
		size := uint32(buf[3]) | uint32(buf[2])<<8 | uint32(buf[1])<<16 | uint32(buf[0])<<24
		buf = buf[4:4+size]
	}
	iter := cfg.NewIterator(nil, buf)
	decoder.Decode(val, iter)
	if iter.Error() != nil {
		return iter.Error()
	}
	return nil
}

func (cfg *frozenConfig) Marshal(obj interface{}) ([]byte, error) {
	msg, isMsg := obj.(protocol.Message)
	if !isMsg {
		return nil, errors.New("can only unmarshal protocol.Message")
	}
	stream := cfg.NewStream(nil, nil)
	stream.WriteMessage(msg)
	if stream.Error() != nil {
		return nil, stream.Error()
	}
	buf := stream.Buffer()
	if cfg.isFramed {
		size := len(buf)
		buf = append([]byte{
			byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
		}, buf...)
	}
	return buf, nil
}

func (cfg *frozenConfig) NewDecoder(reader io.Reader, buf []byte) Decoder {
	if cfg.isFramed {
		return &framedDecoder{reader: reader, iter: cfg.NewIterator(nil, nil)}
	} else {
		return &unframedDecoder{cfg: cfg, iter: cfg.NewIterator(reader, buf), decodeFromReader: reader != nil}
	}
}

func (cfg *frozenConfig) NewEncoder(writer io.Writer) Encoder {
	if cfg.isFramed {
		return &framedEncoder{writer: writer, stream: cfg.NewStream(nil, nil)}
	} else {
		return &unframedEncoder{stream: cfg.NewStream(writer, nil)}
	}
}
