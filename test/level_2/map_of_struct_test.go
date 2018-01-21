package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/test/level_2/map_of_struct_test"
)

func Test_skip_map_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.STRUCT, 1)
		proto.WriteI64(1)

		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()

		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMap(nil))
	}
}

func Test_decode_map_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.STRUCT, 1)
		proto.WriteI64(1)

		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()

		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(map[protocol.FieldId]interface{}{
			protocol.FieldId(1): int64(1024),
		}, iter.ReadMap()[int64(1)])
	}
}

func Test_unmarshal_map_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.STRUCT, 1)
		proto.WriteI64(1)

		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()

		proto.WriteMapEnd()
		var val map[int64]map_of_struct_test.TestObject
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(map[int64]map_of_struct_test.TestObject{
			1: {1024},
		}, val)
	}
}

func Test_encode_map_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteMap(map[interface{}]interface{}{
			int64(1): map[protocol.FieldId]interface{} {
				protocol.FieldId(1): int64(1024),
			},
		})
		iter := c.CreateIterator(stream.Buffer())
		should.Equal(map[protocol.FieldId]interface{}{
			protocol.FieldId(1): int64(1024),
		}, iter.ReadMap()[int64(1)])
	}
}

func Test_marshal_map_of_struct(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal(map[int64]map_of_struct_test.TestObject{
			1: {1024},
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal(map[protocol.FieldId]interface{}{
			protocol.FieldId(1): int64(1024),
		}, iter.ReadMap()[int64(1)])
	}
}