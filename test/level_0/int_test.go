package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_int(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI64(-1)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(int(-1), iter.ReadInt())
	}
}

func Test_unmarshal_int(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI64(-1)
		var val int
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(int(-1), val)
	}
}

func Test_encode_int(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteInt(-1)
		iter := c.CreateIterator(stream.Buffer())
		should.Equal(int(-1), iter.ReadInt())
	}
}

func Test_marshal_int(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal(int(-1))
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal(int(-1), iter.ReadInt())
	}
}
