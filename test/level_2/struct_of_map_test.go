package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/test/level_2/struct_of_map_test"
)

func Test_skip_struct_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.MAP, 1)
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteMapEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipStruct(nil))
	}
}

func Test_decode_struct_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.MAP, 1)
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteMapEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(map[interface{}]interface{}{
			int32(2): int64(2),
		}, iter.ReadStruct()[protocol.FieldId(1)])
	}
}

func Test_unmarshal_struct_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.MAP, 1)
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteMapEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val struct_of_map_test.TestObject
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(struct_of_map_test.TestObject{
			map[int32]int64{2: 2},
		}, val)
	}
}

func Test_encode_struct_of_map(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteStruct(map[protocol.FieldId]interface{}{
		protocol.FieldId(1): map[interface{}]interface{}{
			int32(2): int64(2),
		},
	})
	iter := thrifter.NewIterator(nil,  stream.Buffer())
	should.Equal(map[interface{}]interface{}{
		int32(2): int64(2),
	}, iter.ReadStruct()[protocol.FieldId(1)])
}