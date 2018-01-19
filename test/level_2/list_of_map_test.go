package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
)

func Test_skip_list_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.MAP, 2)
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteMapEnd()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteMapEnd()
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipList(nil))
	}
}

func Test_decode_list_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.MAP, 2)
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteMapEnd()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteMapEnd()
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(map[interface{}]interface{}{
			int32(1): int64(1),
		}, iter.ReadList()[0])
	}
}

func Test_unmarshal_list_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.MAP, 2)
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteMapEnd()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteMapEnd()
		proto.WriteListEnd()
		var val []map[int32]int64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal([]map[int32]int64{
			{1:1}, {2:2},
		}, val)
	}
}

func Test_encode_list_of_map(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteList([]interface{}{
		map[interface{}]interface{} {
			int32(1): int64(1),
		},
		map[interface{}]interface{} {
			int32(2): int64(2),
		},
	})
	iter := thrifter.NewIterator(nil,  stream.Buffer())
	should.Equal(map[interface{}]interface{}{
		int32(1): int64(1),
	}, iter.ReadList()[0])
}