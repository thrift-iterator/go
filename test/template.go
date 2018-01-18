package test

import (
	"github.com/thrift-iterator/go"
	"bytes"
	"git.apache.org/thrift.git/lib/go/thrift"
	"reflect"
)

type Combination struct {
	CreateProtocol func() (*thrift.TMemoryBuffer, thrift.TProtocol)
	CreateIterator func(buf []byte) thrifter.Iterator
	Unmarshal      func(buf []byte, obj interface{}) error
}

var Combinations = []Combination{
	{
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTBinaryProtocol(buf, true, true)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			return thrifter.Config{Protocol: thrifter.ProtocolBinary}.Froze().NewIterator(nil, buf)
		},
		Unmarshal: func(buf []byte, obj interface{}) error {
			cfg := thrifter.Config{Protocol: thrifter.ProtocolBinary}
			cfg = cfg.Decode(reflect.TypeOf(obj))
			return cfg.Froze().Unmarshal(buf, obj)
		},
	},
	{
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTBinaryProtocol(buf, true, true)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			return thrifter.NewIterator(bytes.NewBuffer(buf), nil)
		},
		Unmarshal: func(buf []byte, obj interface{}) error {
			cfg := thrifter.Config{Protocol: thrifter.ProtocolBinary, DecodeFromReader: true}
			cfg = cfg.Decode(reflect.TypeOf(obj))
			api := cfg.Froze()
			decoder := api.NewDecoder(bytes.NewBuffer(buf))
			return decoder.Decode(obj)
		},
	},
	{
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTCompactProtocol(buf)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			cfg := thrifter.Config{Protocol: thrifter.ProtocolCompact}.Froze()
			return cfg.NewIterator(nil, buf)
		},
		Unmarshal: func(buf []byte, obj interface{}) error {
			cfg := thrifter.Config{Protocol: thrifter.ProtocolCompact}
			cfg = cfg.Decode(reflect.TypeOf(obj))
			return cfg.Froze().Unmarshal(buf, obj)
		},
	},
}
