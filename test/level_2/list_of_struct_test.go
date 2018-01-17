package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
)

func Test_skip_list_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRUCT, 2)
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipList(nil))
	}
}

func Test_decode_list_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRUCT, 2)
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(map[protocol.FieldId]interface{}{
			protocol.FieldId(1): int64(1024),
		}, iter.ReadList()[0])
	}
}

func Test_encode_list_of_struct(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteList([]interface{}{
		map[protocol.FieldId]interface{}{
			protocol.FieldId(1): int64(1024),
		},
		map[protocol.FieldId]interface{}{
			protocol.FieldId(1): int64(1024),
		},
	})
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	should.Equal(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): int64(1024),
	}, iter.ReadList()[0])
}
