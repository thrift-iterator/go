package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
	"github.com/v2pro/wombat"
)

func Test_decode_uint16(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI16(1024)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(uint16(1024), iter.ReadUint16())
	}
}

func Test_unmarshal_uint16(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI16(1024)
		var val uint16
		cfg := c.Config.Decode(wombat.Uint16)
		should.NoError(c.Unmarshal(cfg, buf.Bytes(), &val))
		should.Equal(uint16(1024), val)
	}
}

func Test_encode_uint16(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteUInt16(1024)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(uint16(1024), iter.ReadUint16())
}
