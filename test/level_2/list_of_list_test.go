package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
)

func Test_skip_list_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.LIST, 2)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(2)
		proto.WriteListEnd()
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipList(nil))
	}
}

func Test_unmarshal_general_list_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.LIST, 2)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(2)
		proto.WriteListEnd()
		proto.WriteListEnd()
		var val []interface{}
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal([]interface{}{int64(1)}, val[0])
	}
}

func Test_unmarshal_list_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.LIST, 2)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(2)
		proto.WriteListEnd()
		proto.WriteListEnd()
		var val [][]int64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal([][]int64{
			{1}, {2},
		}, val)
	}
}

func Test_marshal_general_list_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal([]interface{}{
			[]interface{}{
				int64(1),
			},
			[]interface{} {
				int64(2),
			},
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal([]interface{}{int64(1)}, iter.ReadList()[0])
	}
}

func Test_marshal_list_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal([][]int64{
			{1}, {2},
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		should.Equal([]interface{}{int64(1)}, iter.ReadList()[0])
	}
}
