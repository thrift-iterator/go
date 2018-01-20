package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/test/level_3/struct_of_list_string"
)

func Test_unmarshal_struct_of_list_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteStructBegin("hello")
		proto.WriteFieldBegin("field1", thrift.LIST, 1)
		proto.WriteListBegin(thrift.STRING, 3)
		proto.WriteString("a")
		proto.WriteString("b")
		proto.WriteString("c")
		proto.WriteListEnd()
		proto.WriteFieldEnd()
		proto.WriteFieldBegin("field2", thrift.I64, 2)
		proto.WriteI64(1024)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		var val struct_of_list_string.TestObject
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(struct_of_list_string.TestObject{
			[]string{"a", "b", "c"},
			1024,
		}, val)
	}
}
