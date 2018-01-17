package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_int64(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI64(-1)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(int64(-1), iter.ReadInt64())
	}
}

func Test_encode_int64(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteInt64(-1)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(int64(-1), iter.ReadInt64())
}
