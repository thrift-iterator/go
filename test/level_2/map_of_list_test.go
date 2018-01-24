package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/general"
)

func Test_skip_map_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.LIST, 1)
		proto.WriteI64(1)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMap(nil))
	}
}

func Test_unmarshal_general_map_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.LIST, 1)
		proto.WriteI64(1)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteMapEnd()
		var val general.Map
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.List{
			int64(1),
		}, val[int64(1)])
	}
}

func Test_unmarshal_map_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.LIST, 1)
		proto.WriteI64(1)
		proto.WriteListBegin(thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteListEnd()
		proto.WriteMapEnd()
		var val map[int64][]int64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(map[int64][]int64{
			1: {1},
		}, val)
	}
}

func Test_marshal_general_map_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(general.Map{
			int64(1): general.List{int64(1)},
		})
		should.NoError(err)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.List{
			int64(1),
		}, val[int64(1)])
	}
}

func Test_marshal_map_of_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal(map[int64][]int64{
			1: {1},
		})
		should.NoError(err)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.List{
			int64(1),
		}, val[int64(1)])
	}
}
