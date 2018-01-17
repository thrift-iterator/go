package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_int8(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteByte(-1)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(int8(-1), iter.ReadInt8())
	}
}

func Test_encode_int8(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteInt8(-1)
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	should.Equal(int8(-1), iter.ReadInt8())
}