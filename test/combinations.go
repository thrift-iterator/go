package test

import (
	"github.com/thrift-iterator/go"
	"bytes"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/spi"
)

type Combination struct {
	CreateProtocol func() (*thrift.TMemoryBuffer, thrift.TProtocol)
	CreateStream   func() spi.Stream
	CreateIterator func(buf []byte) spi.Iterator
	Unmarshal      func(buf []byte, obj interface{}) error
	Marshal        func(obj interface{}) ([]byte, error)
}

var binaryCfg = thrifter.Config{Protocol: thrifter.ProtocolBinary}
var binary = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTBinaryProtocol(buf, true, true)
		return buf, proto
	},
	CreateStream: func() spi.Stream {
		return binaryCfg.Froze().NewStream(nil, nil)
	},
	CreateIterator: func(buf []byte) spi.Iterator {
		return binaryCfg.Froze().NewIterator(nil, buf)
	},
	Unmarshal: func(buf []byte, obj interface{}) error {
		return binaryCfg.Froze().Unmarshal(buf, obj)
	},
	Marshal: func(obj interface{}) ([]byte, error) {
		return binaryCfg.Froze().Marshal(obj)
	},
}

var binaryEncoderDecoder = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTBinaryProtocol(buf, true, true)
		return buf, proto
	},
	CreateStream: func() spi.Stream {
		return binaryCfg.Froze().NewStream(nil, nil)
	},
	CreateIterator: func(buf []byte) spi.Iterator {
		return binaryCfg.Froze().NewIterator(bytes.NewBuffer(buf), nil)
	},
	Unmarshal: func(buf []byte, obj interface{}) error {
		decoder := binaryCfg.Froze().NewDecoder(bytes.NewBuffer(buf), nil)
		return decoder.Decode(obj)
	},
	Marshal: func(obj interface{}) ([]byte, error) {
		encoder := binaryCfg.Froze().NewEncoder(nil)
		err := encoder.Encode(obj)
		if err != nil {
			return nil, err
		}
		return encoder.Buffer(), nil
	},
}

var compactCfg = thrifter.Config{Protocol: thrifter.ProtocolCompact}
var compact = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTCompactProtocol(buf)
		return buf, proto
	},
	CreateStream: func() spi.Stream {
		return compactCfg.Froze().NewStream(nil, nil)
	},
	CreateIterator: func(buf []byte) spi.Iterator {
		return compactCfg.Froze().NewIterator(nil, buf)
	},
	Unmarshal: func(buf []byte, obj interface{}) error {
		return compactCfg.Froze().Unmarshal(buf, obj)
	},
	Marshal: func(obj interface{}) ([]byte, error) {
		return compactCfg.Froze().Marshal(obj)
	},
}
var binaryDynamicCfg = thrifter.Config{Protocol: thrifter.ProtocolBinary, DynamicCodegen: true}
var binaryDynamic = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTBinaryProtocol(buf, true, true)
		return buf, proto
	},
	Unmarshal: func(buf []byte, obj interface{}) error {
		return binaryDynamicCfg.Froze().Unmarshal(buf, obj)
	},
}
var compactDynamicCfg = thrifter.Config{Protocol: thrifter.ProtocolCompact, DynamicCodegen: true}
var compactDynamic = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTCompactProtocol(buf)
		return buf, proto
	},
	Unmarshal: func(buf []byte, obj interface{}) error {
		return compactDynamicCfg.Froze().Unmarshal(buf, obj)
	},
}

var Combinations = []Combination{
	binary, binaryEncoderDecoder, compact,
}

var UnmarshalCombinations = append(Combinations,
	binaryDynamic, compactDynamic)
var MarshalCombinations = UnmarshalCombinations
