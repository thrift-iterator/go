package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
	"github.com/v2pro/wombat"
)

func Test_decode_int32(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI32(-1)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(int32(-1), iter.ReadInt32())
	}
}

func Test_unmarshal_int32(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI32(-1)
		var val int32
		api := c.Config.Decode(wombat.Int32).Froze()
		should.NoError(api.Unmarshal(buf.Bytes(), &val))
		should.Equal(int32(-1), val)
	}
}

func Test_encode_int32(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteInt32(-1)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(int32(-1), iter.ReadInt32())
}