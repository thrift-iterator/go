package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/test/level_2/struct_of_string_test"
)

func Test_skip_struct_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.STRING, 1)
		proto.WriteString("abc")
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipStruct(nil))
	}
}

func Test_unmarshal_general_struct_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.STRING, 1)
		proto.WriteString("abc")
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val map[protocol.FieldId]interface{}
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal("abc", val[protocol.FieldId(1)])
	}
}

func Test_unmarshal_struct_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.STRING, 1)
		proto.WriteString("abc")
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val struct_of_string_test.TestObject
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(struct_of_string_test.TestObject{
			"abc",
		}, val)
	}
}

func Test_marshal_general_struct_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(map[protocol.FieldId]interface{}{
			protocol.FieldId(1): "abc",
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal("abc", iter.ReadStruct()[protocol.FieldId(1)])
	}
}

func Test_marshal_struct_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal(struct_of_string_test.TestObject{
			"abc",
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal("abc", iter.ReadStruct()[protocol.FieldId(1)])
	}
}