package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_uint8(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteByte(100)
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(uint8(100), iter.ReadUInt8())
}

func Test_encode_uint8(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteUInt8(100)
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(uint8(100), iter.ReadUInt8())
}
