package test

import (
	"github.com/thrift-iterator/go"
	"bytes"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type Combination struct {
	CreateProtocol func() (*thrift.TMemoryBuffer, thrift.TProtocol)
	CreateIterator func(buf []byte) thrifter.Iterator
}

var Combinations = []Combination{
	{
		CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
			buf := thrift.NewTMemoryBuffer()
			proto := thrift.NewTBinaryProtocol(buf, true, true)
			return buf, proto
		},
		CreateIterator: func(buf []byte) thrifter.Iterator {
			return thrifter.NewIterator(nil, buf)
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
	},
}
