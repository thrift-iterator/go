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

func Test_marshal_struct_of_1_ptr(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		one := 1
		obj := struct_of_pointer_test.StructOf1Ptr{
			&one,
		}

		output, err := c.Marshal(obj)
		should.NoError(err)
		output1, err := c.Marshal(&obj)
		should.NoError(err)
		should.Equal(output, output1)
		var val *struct_of_pointer_test.StructOf1Ptr
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(1, *val.Field1)
	}
}

func Test_marshal_struct_of_2_ptr(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		one := 1
		two := 2
		obj := struct_of_pointer_test.StructOf2Ptr{
			&one, &two,
		}

		output, err := c.Marshal(obj)
		should.NoError(err)
		output1, err := c.Marshal(&obj)
		should.NoError(err)
		should.Equal(output, output1)
		var val *struct_of_pointer_test.StructOf2Ptr
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(1, *val.Field1)
		should.Equal(2, *val.Field2)
	}
}