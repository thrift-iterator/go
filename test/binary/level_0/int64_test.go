package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_int64(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteI64(-1)
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(int64(-1), iter.ReadInt64())
}

func Test_encode_int64(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteInt64(-1)
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(int64(-1), iter.ReadInt64())
}
