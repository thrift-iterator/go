package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
)

func Test_skip_struct_of_map(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.MAP, 1)
	proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
	proto.WriteI32(2)
	proto.WriteI64(2)
	proto.WriteMapEnd()
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipStruct(nil))
}

func Test_decode_struct_of_map(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.MAP, 1)
	proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
	proto.WriteI32(2)
	proto.WriteI64(2)
	proto.WriteMapEnd()
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	should.Equal(map[interface{}]interface{}{
		int32(2): int64(2),
	}, iter.ReadStruct()[protocol.FieldId(1)])
}

func Test_encode_struct_of_map(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteStruct(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): map[interface{}]interface{}{
			int32(2): int64(2),
		},
	})
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	should.Equal(map[interface{}]interface{}{
		int32(2): int64(2),
	}, iter.ReadStruct()[protocol.FieldId(1)])
}