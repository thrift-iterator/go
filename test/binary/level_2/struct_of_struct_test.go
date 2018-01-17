package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
)

func Test_skip_struct_of_struct(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.STRUCT, 1)

	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.STRING, 1)
	proto.WriteString("abc")
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()

	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipStruct())
}

func Test_decode_struct_of_struct(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.STRUCT, 1)

	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.STRING, 1)
	proto.WriteString("abc")
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()

	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): "abc",
	}, iter.ReadStruct()[protocol.FieldId(1)])
}

func Test_encode_struct_of_struct(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteStruct(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): map[protocol.FieldId]interface{}{
			protocol.FieldId(1): "abc",
		},
	})
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): "abc",
	}, iter.ReadStruct()[protocol.FieldId(1)])
}