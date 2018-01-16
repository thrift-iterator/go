package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_int32(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteI32(-1)
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(int32(-1), iter.ReadInt32())
}

func Test_encode_int32(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteInt32(-1)
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(int32(-1), iter.ReadInt32())
}