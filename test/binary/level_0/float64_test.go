package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_float64(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteDouble(10.24)
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(10.24, iter.ReadFloat64())
}

func Test_encode_float64(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteFloat64(10.24)
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(10.24, iter.ReadFloat64())
}