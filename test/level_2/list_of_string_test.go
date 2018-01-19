package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
)

func Test_skip_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRING, 3)
		proto.WriteString("a")
		proto.WriteString("b")
		proto.WriteString("c")
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipList(nil))
	}
}

func Test_decode_list_of_string(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.STRING, 3)
	proto.WriteString("a")
	proto.WriteString("b")
	proto.WriteString("c")
	proto.WriteListEnd()
	iter := thrifter.NewIterator(nil,  buf.Bytes())
	should.Equal([]interface{}{"a", "b", "c"}, iter.ReadList())
}

func Test_unmarshal_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRING, 3)
		proto.WriteString("a")
		proto.WriteString("b")
		proto.WriteString("c")
		proto.WriteListEnd()
		var val []string
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal([]string{
			"a", "b", "c",
		}, val)
	}
}

func Test_encode_list_of_string(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteList([]interface{}{
		"a", "b", "c",
	})
	iter := thrifter.NewIterator(nil,  stream.Buffer())
	should.Equal([]interface{}{"a", "b", "c"}, iter.ReadList())
}
