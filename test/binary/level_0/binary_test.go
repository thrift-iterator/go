package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_binary(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteBinary([]byte("hello"))
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal("hello", string(iter.ReadBinary()))
}

func Test_encode_binary(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteBinary([]byte("hello"))
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal("hello", string(iter.ReadBinary()))
}
