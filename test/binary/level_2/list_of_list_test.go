package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_skip_list_of_list(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.LIST, 2)
	proto.WriteListBegin(thrift.I64, 1)
	proto.WriteI64(1)
	proto.WriteListEnd()
	proto.WriteListBegin(thrift.I64, 1)
	proto.WriteI64(2)
	proto.WriteListEnd()
	proto.WriteListEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipList())
}

func Test_decode_list_of_list(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.LIST, 2)
	proto.WriteListBegin(thrift.I64, 1)
	proto.WriteI64(1)
	proto.WriteListEnd()
	proto.WriteListBegin(thrift.I64, 1)
	proto.WriteI64(2)
	proto.WriteListEnd()
	proto.WriteListEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal([]interface{}{int64(1)}, iter.ReadList()[0])
}

func Test_encode_list_of_list(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteList([]interface{}{
		[]interface{}{
			int64(1),
		},
		[]interface{} {
			int64(2),
		},
	})
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal([]interface{}{int64(1)}, iter.ReadList()[0])
}
