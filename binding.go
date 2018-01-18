package thrifter

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go/binding"
	"github.com/thrift-iterator/go/protocol/binary"
	"github.com/thrift-iterator/go/protocol/compact"
	"github.com/thrift-iterator/go/protocol/sbinary"
)

type funcDecoder struct {
	f func(dst interface{}, src interface{})
}

func (decoder *funcDecoder) Decode(val interface{}, iter Iterator) {
	decoder.f(val, iter)
}

func (cfg Config) Decode(typ reflect.Type) Config {
	iteratorType := reflect.TypeOf((*binary.Iterator)(nil))
	if cfg.DecodeFromReader {
		iteratorType = reflect.TypeOf((*sbinary.Iterator)(nil))
	}
	if cfg.Protocol == ProtocolCompact {
		iteratorType = reflect.TypeOf((*compact.Iterator)(nil))
	}
	funcObj := generic.Expand(binding.Decode,
		"ST", iteratorType,
		"DT", typ)
	f := funcObj.(func(interface{}, interface{}))
	if cfg.Decoders == nil {
		cfg.Decoders = map[reflect.Type]ValDecoder{}
	}
	cfg.Decoders[typ] = &funcDecoder{f}
	return cfg
}
