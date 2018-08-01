package test

import (
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/test"
	"testing"
)

func Test_decode_bool(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteBool(true)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(true, iter.ReadBool())

		buf, proto = c.CreateProtocol()
		proto.WriteBool(false)
		iter = c.CreateIterator(buf.Bytes())
		should.Equal(false, iter.ReadBool())
	}
}

func Test_unmarshal_bool(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		var val1 bool
		proto.WriteBool(true)
		should.NoError(c.Unmarshal(buf.Bytes(), &val1))
		should.Equal(true, val1)

		buf, proto = c.CreateProtocol()
		var val2 bool = true
		proto.WriteBool(false)
		should.NoError(c.Unmarshal(buf.Bytes(), &val2))
		should.Equal(false, val2)
	}
}

func Test_encode_bool(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteBool(true)
		iter := c.CreateIterator(stream.Buffer())
		should.Equal(true, iter.ReadBool())

		stream = c.CreateStream()
		stream.WriteBool(false)
		iter = c.CreateIterator(stream.Buffer())
		should.Equal(false, iter.ReadBool())
	}
}

func Test_marshal_bool(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal(true)
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal(true, iter.ReadBool())

		output, err = c.Marshal(false)
		should.NoError(err)
		iter = c.CreateIterator(output)
		should.Equal(false, iter.ReadBool())
	}
}
