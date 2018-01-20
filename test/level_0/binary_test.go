package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_binary(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteBinary([]byte("hello"))
		iter := c.CreateIterator(buf.Bytes())
		should.Equal("hello", string(iter.ReadBinary()))
	}
}

func Test_unmarshal_binary(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteBinary([]byte("hello"))
		var val []byte
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal("hello", string(val))
	}
}

func Test_encode_binary(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteBinary([]byte(`hello world!`))
		iter := c.CreateIterator(stream.Buffer())
		should.Equal([]byte(`hello world!`), iter.ReadBinary())
	}
}

func Test_marshal_binary(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		val := []byte("hello")
		output, err := c.Marshal(val)
		should.NoError(err)
		iter := thrifter.NewIterator(nil, output)
		should.Equal("hello", string(iter.ReadBinary()))
	}
}
