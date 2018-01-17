package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_uint8(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteByte(100)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(uint8(100), iter.ReadUInt8())
	}
}

func Test_encode_uint8(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteUInt8(100)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(uint8(100), iter.ReadUInt8())
}
