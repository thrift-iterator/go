package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_skip_list_of_map(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.MAP, 2)
	proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
	proto.WriteI32(1)
	proto.WriteI64(1)
	proto.WriteMapEnd()
	proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
	proto.WriteI32(2)
	proto.WriteI64(2)
	proto.WriteMapEnd()
	proto.WriteListEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipList())
}

func Test_decode_list_of_map(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.MAP, 2)
	proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
	proto.WriteI32(1)
	proto.WriteI64(1)
	proto.WriteMapEnd()
	proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
	proto.WriteI32(2)
	proto.WriteI64(2)
	proto.WriteMapEnd()
	proto.WriteListEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(map[interface{}]interface{}{
		int32(1): int64(1),
	}, iter.ReadList()[0])
}

func Test_encode_list_of_map(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteList([]interface{}{
		map[interface{}]interface{} {
			int32(1): int64(1),
		},
		map[interface{}]interface{} {
			int32(2): int64(2),
		},
	})
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(map[interface{}]interface{}{
		int32(1): int64(1),
	}, iter.ReadList()[0])
}