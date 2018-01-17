package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_bool(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteBool(true)
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(true, iter.ReadBool())
	}
}

func Test_encode_bool(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteBool(true)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	should.Equal(true, iter.ReadBool())
}
