package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_binary(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteBinary([]byte("hello"))
		iter := c.CreateIterator(buf.Bytes())
		should.Equal("hello", string(iter.ReadBinary()))
	}
}

func Test_encode_binary(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteBinary([]byte("hello"))
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal("hello", string(iter.ReadBinary()))
}
