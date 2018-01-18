package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_int16(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI16(-1)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(int16(-1), iter.ReadInt16())
	}
}

func Test_unmarshal_int16(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI16(-1)
		var val int16
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(int16(-1), val)
	}
}

func Test_encode_int16(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteInt16(-1)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(int16(-1), iter.ReadInt16())
}