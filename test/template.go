package test

import (
	"github.com/thrift-iterator/go"
	"bytes"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type Combination struct {
	CreateProtocol func() (*thrift.TMemoryBuffer, thrift.TProtocol)
	CreateIterator func(buf []byte) thrifter.Iterator
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
	},
	{
		Config: thrifter.Config{Protocol: thrifter.ProtocolBinary},
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTBinaryProtocol(buf, true, true)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			return thrifter.NewIterator(bytes.NewBuffer(buf), nil)
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
	},
}
