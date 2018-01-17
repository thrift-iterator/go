package test

import (
	"testing"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/protocol"
)

func Test_decode_struct_by_iterator(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.I64, 1)
	proto.WriteI64(1024)
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	called := false
	iter.ReadStructCB(func(fieldType protocol.TType, fieldId protocol.FieldId) {
		should.False(called)
		called = true
		should.Equal(protocol.I64, fieldType)
		should.Equal(protocol.FieldId(1), fieldId)
		should.Equal(int64(1024), iter.ReadInt64())
	})
	should.True(called)
}

func Test_encode_struct_by_stream(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteStructField(protocol.I64, protocol.FieldId(1))
	stream.WriteInt64(1024)
	stream.WriteStructFieldStop()
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	called := false
	iter.ReadStructCB(func(fieldType protocol.TType, fieldId protocol.FieldId) {
		should.False(called)
		called = true
		should.Equal(protocol.I64, fieldType)
		should.Equal(protocol.FieldId(1), fieldId)
		should.Equal(int64(1024), iter.ReadInt64())
	})
}

func Test_decode_struct_as_object(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.I64, 1)
	proto.WriteI64(1024)
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	obj := iter.ReadStruct()
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): int64(1024),
	}, obj)
}

func Test_encode_struct_from_object(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteStruct(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): int64(1024),
	})
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	obj := iter.ReadStruct()
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): int64(1024),
	}, obj)
}

func Test_skip_struct(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.I64, 1)
	proto.WriteI64(1024)
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipStruct())
}
