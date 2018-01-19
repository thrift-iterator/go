package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
)

func Test_decode_list_by_iterator(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.I64, 3)
		proto.WriteI64(1)
		proto.WriteI64(2)
		proto.WriteI64(3)
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		elemType, length := iter.ReadListHeader()
		should.Equal(protocol.TypeI64, elemType)
		should.Equal(3, length)
		should.Equal(uint64(1), iter.ReadUint64())
		should.Equal(uint64(2), iter.ReadUint64())
		should.Equal(uint64(3), iter.ReadUint64())
	}
}

func Test_encode_list_by_stream(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteListHeader(protocol.TypeI64, 3)
	stream.WriteUInt64(1)
	stream.WriteUInt64(2)
	stream.WriteUInt64(3)
	iter := thrifter.NewIterator(nil, stream.Buffer())
	elemType, length := iter.ReadListHeader()
	should.Equal(protocol.TypeI64, elemType)
	should.Equal(3, length)
	should.Equal(uint64(1), iter.ReadUint64())
	should.Equal(uint64(2), iter.ReadUint64())
	should.Equal(uint64(3), iter.ReadUint64())
}

func Test_decode_list_as_object(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.I64, 3)
		proto.WriteI64(1)
		proto.WriteI64(2)
		proto.WriteI64(3)
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		obj := iter.ReadList()
		should.Equal([]interface{}{int64(1), int64(2), int64(3)}, obj)
	}
}

func Test_unmarshal_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.I64, 3)
		proto.WriteI64(1)
		proto.WriteI64(2)
		proto.WriteI64(3)
		proto.WriteListEnd()
		var val []int64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal([]int64{int64(1), int64(2), int64(3)}, val)
	}
}

func Test_encode_list_from_object(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewStream(nil, nil)
	stream.WriteList([]interface{}{
		int64(1), int64(2), int64(3),
	})
	iter := thrifter.NewIterator(nil, stream.Buffer())
	elemType, length := iter.ReadListHeader()
	should.Equal(protocol.TypeI64, elemType)
	should.Equal(3, length)
	should.Equal(uint64(1), iter.ReadUint64())
	should.Equal(uint64(2), iter.ReadUint64())
	should.Equal(uint64(3), iter.ReadUint64())
}

func Test_skip_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.I64, 3)
		proto.WriteI64(1)
		proto.WriteI64(2)
		proto.WriteI64(3)
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipList(nil))
	}
}
