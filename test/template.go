package test

import (
	"github.com/thrift-iterator/go"
	"bytes"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type Combination struct {
	CreateProtocol func() (*thrift.TMemoryBuffer, thrift.TProtocol)
	CreateIterator func(buf []byte) thrifter.Iterator
	Unmarshal      func(cfg thrifter.Config, buf []byte, obj interface{}) error
	Config         thrifter.Config
}

var Combinations = []Combination{
	{
		Config: thrifter.Config{Protocol: thrifter.ProtocolBinary},
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTBinaryProtocol(buf, true, true)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			return thrifter.Config{Protocol: thrifter.ProtocolBinary}.Froze().NewIterator(nil, buf)
		},
		Unmarshal: func(cfg thrifter.Config, buf []byte, obj interface{}) error {
			return cfg.Froze().Unmarshal(buf, obj)
		},
	},
	{
		Config: thrifter.Config{Protocol: thrifter.ProtocolBinary, DecodeFromReader: true},
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTBinaryProtocol(buf, true, true)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			return thrifter.NewIterator(bytes.NewBuffer(buf), nil)
		},
		Unmarshal: func(cfg thrifter.Config, buf []byte, obj interface{}) error {
			api := cfg.Froze()
			decoder := api.NewDecoder(bytes.NewBuffer(buf))
			return decoder.Decode(obj)
		},
	},
	{
		Config: thrifter.Config{Protocol: thrifter.ProtocolCompact},
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTCompactProtocol(buf)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			cfg := thrifter.Config{Protocol: thrifter.ProtocolCompact}.Froze()
			return cfg.NewIterator(nil, buf)
		},
		Unmarshal: func(cfg thrifter.Config, buf []byte, obj interface{}) error {
			return cfg.Froze().Unmarshal(buf, obj)
		},
	},
}
