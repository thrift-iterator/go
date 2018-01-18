package thrifter

import (
	"reflect"
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go/binding"
	"github.com/thrift-iterator/go/protocol/binary"
)

type funcDecoder struct {
	f func(dst interface{}, src interface{})
}

func (decoder *funcDecoder) Decode(val interface{}, iter Iterator) {
	decoder.f(val, iter)
}

func (cfg Config) Decode(typ reflect.Type) Config {
	typ = reflect.PtrTo(typ)
	funcObj := generic.Expand(binding.DecodeAnything,
		"ST", reflect.TypeOf((*binary.Iterator)(nil)),
		"DT", typ)
	f := funcObj.(func(interface{}, interface{}))
	if cfg.Decoders == nil {
		cfg.Decoders = map[reflect.Type]ValDecoder{}
	}
	cfg.Decoders[typ] = &funcDecoder{f}
	return cfg
}
