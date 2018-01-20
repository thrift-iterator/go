package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
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
	stream := thrifter.NewStream(nil, nil)
	stream.WriteUInt(1024)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(uint(1024), iter.ReadUint())
}
