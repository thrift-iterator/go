package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_uint32(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI32(1024)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(uint32(1024), iter.ReadUInt32())
	}
}

func Test_encode_uint32(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteUInt32(1024)
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	should.Equal(uint32(1024), iter.ReadUInt32())
}
