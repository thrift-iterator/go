package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_string(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteString("hello")
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal("hello", iter.ReadString())
}

func Test_encode_string(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteString("hello")
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal("hello", iter.ReadString())
}