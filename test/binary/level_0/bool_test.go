package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_bool(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteBool(true)
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(true, iter.ReadBool())
}

func Test_encode_bool(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteBool(true)
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(true, iter.ReadBool())
}