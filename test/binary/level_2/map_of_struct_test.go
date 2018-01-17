package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
)

func Test_skip_map_of_struct(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteMapBegin(thrift.I64, thrift.STRUCT, 1)
	proto.WriteI64(1)

	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.I64, 1)
	proto.WriteI64(1024)
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()

	proto.WriteMapEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipMap())
}

func Test_decode_map_of_struct(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteMapBegin(thrift.I64, thrift.STRUCT, 1)
	proto.WriteI64(1)

	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.I64, 1)
	proto.WriteI64(1024)
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()

	proto.WriteMapEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): int64(1024),
	}, iter.ReadMap()[int64(1)])
}

func Test_encode_map_of_struct(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil)
	stream.WriteMap(map[interface{}]interface{}{
		int64(1): map[protocol.FieldId]interface{} {
			protocol.FieldId(1): int64(1024),
		},
	})
	iter := thrifter.NewIterator(stream.Buffer())
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): int64(1024),
	}, iter.ReadMap()[int64(1)])
}