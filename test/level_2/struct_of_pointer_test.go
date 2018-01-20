package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/test"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test/level_2/struct_of_pointer_test"
)

func Test_unmarshal_struct_of_1_ptr(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val *struct_of_pointer_test.StructOf1Ptr
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(1, *val.Field1)
	}
}

func Test_unmarshal_struct_of_2_ptr(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteFieldEnd()
		proto.WriteFieldBegin("field2", thrift.I64, 2)
		proto.WriteI64(2)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val *struct_of_pointer_test.StructOf2Ptr
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(1, *val.Field1)
		should.Equal(2, *val.Field2)
	}
}