package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteString("hello")
		iter := c.CreateIterator(buf.Bytes())
		should.Equal("hello", iter.ReadString())
	}
}

func Test_encode_string(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteString("hello")
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	should.Equal("hello", iter.ReadString())
}
