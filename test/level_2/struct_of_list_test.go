package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/test/level_2/struct_of_list_test"
	"github.com/thrift-iterator/go/general"
)

func Test_skip_struct_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.LIST, 1)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipStruct(nil))
	}
}

func Test_unmarshal_general_struct_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.LIST, 1)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val general.Struct
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.List{int64(1)}, val[protocol.FieldId(1)])
	}
}

func Test_unmarshal_struct_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.LIST, 1)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val struct_of_list_test.TestObject
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(struct_of_list_test.TestObject{
			[]int64{1},
		}, val)
	}
}

func Test_marshal_general_struct_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		obj := general.Struct {
			protocol.FieldId(1): general.List {
				int64(1),
			},
		}

		output, err := c.Marshal(obj)
		should.NoError(err)
		output1, err := c.Marshal(&obj)
		should.NoError(err)
		should.Equal(output, output1)
		var val general.Struct
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.List{int64(1)}, val[protocol.FieldId(1)])
	}
}

func Test_marshal_struct_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		obj := struct_of_list_test.TestObject{
			[]int64{1},
		}

		output, err := c.Marshal(obj)
		should.NoError(err)
		output1, err := c.Marshal(&obj)
		should.NoError(err)
		should.Equal(output, output1)
		var val general.Struct
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.List{int64(1)}, val[protocol.FieldId(1)])
	}
}
