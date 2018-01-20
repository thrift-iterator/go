package test

import (
	"testing"
	"github.com/stretchr/testify/require"
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

func Test_unmarshal_float64(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteDouble(10.24)
		var val float64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(10.24, val)
	}
}

func Test_encode_float64(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteFloat64(10.24)
		iter := c.CreateIterator(stream.Buffer())
		should.Equal(10.24, iter.ReadFloat64())
	}
}
