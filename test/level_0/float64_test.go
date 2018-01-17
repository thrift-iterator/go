package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_float64(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteDouble(10.24)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(10.24, iter.ReadFloat64())
	}
}

func Test_encode_float64(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteFloat64(10.24)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(10.24, iter.ReadFloat64())
}
