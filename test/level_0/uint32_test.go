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
		should.Equal(uint32(1024), iter.ReadUint32())
	}
}

func Test_unmarshal_uint32(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI32(1024)
		var val uint32
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(uint32(1024), val)
	}
}

func Test_encode_uint32(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteUInt32(1024)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(uint32(1024), iter.ReadUint32())
}
