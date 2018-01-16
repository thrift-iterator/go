package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_uint64(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteI64(1024)
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(uint64(1024), iter.ReadUInt64())
}

func Test_encode_uint64(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteUInt64(1024)
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(uint64(1024), iter.ReadUInt64())
}
