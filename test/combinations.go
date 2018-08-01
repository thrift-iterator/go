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
	Unmarshal      func(buf []byte, val interface{}) error
	Marshal        func(val interface{}) ([]byte, error)
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
	Unmarshal: func(buf []byte, val interface{}) error {
		return binaryCfg.Froze().Unmarshal(buf, val)
	},
	Marshal: func(val interface{}) ([]byte, error) {
		return binaryCfg.Froze().Marshal(val)
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
	Unmarshal: func(buf []byte, val interface{}) error {
		decoder := binaryCfg.Froze().NewDecoder(bytes.NewBuffer(buf), nil)
		return decoder.Decode(val)
	},
	Marshal: func(val interface{}) ([]byte, error) {
		encoder := binaryCfg.Froze().NewEncoder(nil)
		err := encoder.Encode(val)
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
	Unmarshal: func(buf []byte, val interface{}) error {
		return compactCfg.Froze().Unmarshal(buf, val)
	},
	Marshal: func(val interface{}) ([]byte, error) {
		return compactCfg.Froze().Marshal(val)
	},
}

var compactEncoderDecoder = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTCompactProtocol(buf)
		return buf, proto
	},
	CreateStream: func() spi.Stream {
		return compactCfg.Froze().NewStream(nil, nil)
	},
	CreateIterator: func(buf []byte) spi.Iterator {
		return compactCfg.Froze().NewIterator(bytes.NewBuffer(buf), nil)
	},
	Unmarshal: func(buf []byte, val interface{}) error {
		decoder := compactCfg.Froze().NewDecoder(bytes.NewBuffer(buf), nil)
		return decoder.Decode(val)
	},
	Marshal: func(val interface{}) ([]byte, error) {
		encoder := compactCfg.Froze().NewEncoder(nil)
		err := encoder.Encode(val)
		if err != nil {
			return nil, err
		}
		return encoder.Buffer(), nil
	},
}

var binaryDynamicCfg = thrifter.Config{Protocol: thrifter.ProtocolBinary, StaticCodegen: false}
var binaryDynamic = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTBinaryProtocol(buf, true, true)
		return buf, proto
	},
	CreateIterator: func(buf []byte) spi.Iterator {
		return binaryDynamicCfg.Froze().NewIterator(nil, buf)
	},
	Unmarshal: func(buf []byte, val interface{}) error {
		return binaryDynamicCfg.Froze().Unmarshal(buf, val)
	},
	Marshal: func(val interface{}) ([]byte, error) {
		return binaryDynamicCfg.Froze().Marshal(val)
	},
}
var compactDynamicCfg = thrifter.Config{Protocol: thrifter.ProtocolCompact, StaticCodegen: false}
var compactDynamic = Combination{
	CreateProtocol: func() (*thrift.TMemoryBuffer, thrift.TProtocol) {
		buf := thrift.NewTMemoryBuffer()
		proto := thrift.NewTCompactProtocol(buf)
		return buf, proto
	},
	CreateIterator: func(buf []byte) spi.Iterator {
		return compactDynamicCfg.Froze().NewIterator(nil, buf)
	},
	Unmarshal: func(buf []byte, val interface{}) error {
		return compactDynamicCfg.Froze().Unmarshal(buf, val)
	},
	Marshal: func(val interface{}) ([]byte, error) {
		return compactDynamicCfg.Froze().Marshal(val)
	},
}

var Combinations = []Combination{
	binary, binaryEncoderDecoder, compact, compactEncoderDecoder,
}

var UnmarshalCombinations = append(Combinations,
	binaryDynamic, compactDynamic)
var MarshalCombinations = UnmarshalCombinations
