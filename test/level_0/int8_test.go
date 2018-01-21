package test

import (
	"testing"
	"github.com/stretchr/testify/require"
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

func Test_unmarshal_int8(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteByte(-1)
		var val int8
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(int8(-1), val)
	}
}

func Test_encode_int8(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteInt8(-1)
		iter := c.CreateIterator(stream.Buffer())
		should.Equal(int8(-1), iter.ReadInt8())
	}
}

func Test_marshal_int8(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(int8(-1))
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal(int8(-1), iter.ReadInt8())
	}
}