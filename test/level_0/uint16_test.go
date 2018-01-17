package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_uint16(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI16(1024)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(uint16(1024), iter.ReadUInt16())
	}
}

func Test_encode_uint16(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteUInt16(1024)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(uint16(1024), iter.ReadUInt16())
}
