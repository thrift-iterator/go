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
	"github.com/thrift-iterator/go/binding/reflection"
	"github.com/thrift-iterator/go/binding/codegen"
	"github.com/thrift-iterator/go/general"
)

type frozenConfig struct {
	extension      spi.Extension
	protocol       Protocol
	genDecoders    unsafe.Pointer
	genEncoders    unsafe.Pointer
	extDecoders    unsafe.Pointer
	extEncoders    unsafe.Pointer
	isFramed       bool
	dynamicCodegen bool
}

func (cfg Config) Froze() API {
	api := &frozenConfig{
		extension:      &general.Extension{},
		protocol:       cfg.Protocol,
		isFramed:       cfg.IsFramed,
		dynamicCodegen: cfg.DynamicCodegen,
	}
	atomic.StorePointer(&api.extDecoders, unsafe.Pointer(&map[string]spi.ValDecoder{}))
	atomic.StorePointer(&api.extEncoders, unsafe.Pointer(&map[string]spi.ValEncoder{}))
	atomic.StorePointer(&api.genDecoders, unsafe.Pointer(&map[reflect.Type]spi.ValDecoder{}))
	atomic.StorePointer(&api.genEncoders, unsafe.Pointer(&map[reflect.Type]spi.ValEncoder{}))
	return api
}

func (cfg *frozenConfig) addGenDecoder(cacheKey reflect.Type, decoder spi.ValDecoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.genDecoders)
		cache := *(*map[reflect.Type]spi.ValDecoder)(ptr)
		copied := map[reflect.Type]spi.ValDecoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = decoder
		done = atomic.CompareAndSwapPointer(&cfg.genDecoders, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) addExtDecoder(cacheKey string, decoder spi.ValDecoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.extDecoders)
		cache := *(*map[string]spi.ValDecoder)(ptr)
		copied := map[string]spi.ValDecoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = decoder
		done = atomic.CompareAndSwapPointer(&cfg.extDecoders, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) addGenEncoder(cacheKey reflect.Type, encoder spi.ValEncoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.genEncoders)
		cache := *(*map[reflect.Type]spi.ValEncoder)(ptr)
		copied := map[reflect.Type]spi.ValEncoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = encoder
		done = atomic.CompareAndSwapPointer(&cfg.genEncoders, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) addExtEncoder(cacheKey string, encoder spi.ValEncoder) {
	done := false
	for !done {
		ptr := atomic.LoadPointer(&cfg.extEncoders)
		cache := *(*map[string]spi.ValEncoder)(ptr)
		copied := map[string]spi.ValEncoder{}
		for k, v := range cache {
			copied[k] = v
		}
		copied[cacheKey] = encoder
		done = atomic.CompareAndSwapPointer(&cfg.extEncoders, ptr, unsafe.Pointer(&copied))
	}
}

func (cfg *frozenConfig) PrepareDecoder(valType reflect.Type) {
	cacheKey := valType.String()
	if cfg.GetDecoder(cacheKey) != nil {
		return
	}
	decoder := cfg.extension.DecoderOf(valType)
	cfg.addExtDecoder(cacheKey, decoder)
	cfg.addGenDecoder(valType, decoder)
}

func (cfg *frozenConfig) GetDecoder(cacheKey string) spi.ValDecoder {
	ptr := atomic.LoadPointer(&cfg.extDecoders)
	cache := *(*map[string]spi.ValDecoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) getGenDecoder(cacheKey reflect.Type) spi.ValDecoder {
	ptr := atomic.LoadPointer(&cfg.genDecoders)
	cache := *(*map[reflect.Type]spi.ValDecoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) PrepareEncoder(valType reflect.Type) {
	cacheKey := valType.String()
	if cfg.GetEncoder(cacheKey) != nil {
		return
	}
	encoder := cfg.extension.EncoderOf(valType)
	cfg.addExtEncoder(cacheKey, encoder)
	cfg.addGenEncoder(valType, encoder)
}

func (cfg *frozenConfig) GetEncoder(cacheKey string) spi.ValEncoder {
	ptr := atomic.LoadPointer(&cfg.extEncoders)
	cache := *(*map[string]spi.ValEncoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) getGenEncoder(cacheKey reflect.Type) spi.ValEncoder {
	ptr := atomic.LoadPointer(&cfg.genEncoders)
	cache := *(*map[reflect.Type]spi.ValEncoder)(ptr)
	return cache[cacheKey]
}

func (cfg *frozenConfig) NewStream(writer io.Writer, buf []byte) spi.Stream {
	switch cfg.protocol {
	case ProtocolBinary:
		return binary.NewStream(cfg, writer, buf)
	case ProtocolCompact:
		return compact.NewStream(cfg, writer, buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) NewIterator(reader io.Reader, buf []byte) spi.Iterator {
	switch cfg.protocol {
	case ProtocolBinary:
		if reader != nil {
			return sbinary.NewIterator(cfg, reader, buf)
		}
		return binary.NewIterator(cfg, buf)
	case ProtocolCompact:
		return compact.NewIterator(cfg, buf)
	}
	panic("unsupported protocol")
}

func (cfg *frozenConfig) WillDecodeFromBuffer(samples ...interface{}) {
	if cfg.dynamicCodegen {
		panic("this config is using dynamic codegen, can not do static codegen")
	}
	for _, sample := range samples {
		cfg.staticDecoderOf(false, reflect.TypeOf(sample))
	}
}

func (cfg *frozenConfig) WillDecodeFromReader(samples ...interface{}) {
	if cfg.dynamicCodegen {
		panic("this config is using dynamic codegen, can not do static codegen")
	}
	for _, sample := range samples {
		cfg.staticDecoderOf(true, reflect.TypeOf(sample))
	}
}

func (cfg *frozenConfig) WillEncode(samples ...interface{}) {
	if cfg.dynamicCodegen {
		panic("this config is using dynamic codegen, can not do static codegen")
	}
	for _, sample := range samples {
		cfg.staticEncoderOf(reflect.TypeOf(sample))
	}
}

func (cfg *frozenConfig) decoderOf(decodeFromReader bool, valType reflect.Type) spi.ValDecoder {
	switch valType {
	case reflect.TypeOf((*map[protocol.FieldId]RawMessage)(nil)):
		return rawStructDecoderInstance
	}
	if cfg.dynamicCodegen {
		return reflection.DecoderOf(cfg.extension, valType)
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
	funcObj := generic.Expand(codegen.Decode,
		"EXT", &codegen.Extension{Extension: cfg.extension},
		"ST", iteratorType,
		"DT", valType)
	f := funcObj.(func(interface{}, interface{}))
	return &funcDecoder{f}
}

func (cfg *frozenConfig) encoderOf(valType reflect.Type) spi.ValEncoder {
	switch valType {
	case reflect.TypeOf((*map[protocol.FieldId]RawMessage)(nil)).Elem():
		return rawStructEncoderInstance
	}
	if cfg.dynamicCodegen {
		return reflection.EncoderOf(cfg.extension, valType)
	}
	return cfg.staticEncoderOf(valType)
}

func (cfg *frozenConfig) staticEncoderOf(valType reflect.Type) spi.ValEncoder {
	streamType := reflect.TypeOf((*binary.Stream)(nil))
	if cfg.protocol == ProtocolCompact {
		streamType = reflect.TypeOf((*compact.Stream)(nil))
	}
	funcObj := generic.Expand(codegen.Encode,
		"EXT", &codegen.Extension{Extension: cfg.extension},
		"ST", valType,
		"DT", streamType)
	f := funcObj.(func(interface{}, interface{}))
	return &funcEncoder{f}
}

type funcDecoder struct {
	f func(dst interface{}, src interface{})
}

func (decoder *funcDecoder) Decode(val interface{}, iter spi.Iterator) {
	decoder.f(val, iter)
}

type funcEncoder struct {
	f func(dst interface{}, src interface{})
}

func (encoder *funcEncoder) Encode(val interface{}, stream spi.Stream) {
	encoder.f(stream, val)
}

func (encoder *funcEncoder) ThriftType() protocol.TType {
	panic("funcEncoder is not composable")
}

func (cfg *frozenConfig) Unmarshal(buf []byte, val interface{}) error {
	valType := reflect.TypeOf(val)
	decoder := cfg.getGenDecoder(valType)
	if decoder == nil {
		decoder = cfg.decoderOf(false, valType)
		cfg.addGenDecoder(valType, decoder)
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

func (cfg *frozenConfig) Marshal(val interface{}) ([]byte, error) {
	valType := reflect.TypeOf(val)
	encoder := cfg.getGenEncoder(valType)
	if encoder == nil {
		encoder = cfg.encoderOf(valType)
		cfg.addGenEncoder(valType, encoder)
	}
	stream := cfg.NewStream(nil, nil)
	encoder.Encode(val, stream)
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
		return &framedDecoder{
			cfg:               cfg,
			shouldDecodeFrame: true,
			reader:            reader,
			iter:              cfg.NewIterator(nil, nil),
		}
	} else {
		return &unframedDecoder{
			cfg:              cfg,
			iter:             cfg.NewIterator(reader, buf),
			decodeFromReader: reader != nil,
		}
	}
}

func (cfg *frozenConfig) NewEncoder(writer io.Writer) Encoder {
	if cfg.isFramed {
		return &framedEncoder{
			cfg:    cfg,
			writer: writer,
			stream: cfg.NewStream(nil, nil),
		}
	} else {
		return &unframedEncoder{
			cfg:    cfg,
			stream: cfg.NewStream(writer, nil),
		}
	}
}
