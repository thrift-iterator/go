package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_uint(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI64(1024)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(uint(1024), iter.ReadUint())
	}
}

func Test_unmarshal_uint(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI64(1024)
		var val uint
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(uint(1024), val)
	}
}

func Test_encode_uint(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteUint(1024)
		iter := c.CreateIterator(stream.Buffer())
		should.Equal(uint(1024), iter.ReadUint())
	}
}

func Test_marshal_uint(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(uint(1024))
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal(uint(1024), iter.ReadUint())
	}
}
