package test

import (
	"testing"
	"github.com/stretchr/testify/require"
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
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteListHeader(protocol.TypeI64, 3)
		stream.WriteUint64(1)
		stream.WriteUint64(2)
		stream.WriteUint64(3)
		iter := c.CreateIterator(stream.Buffer())
		elemType, length := iter.ReadListHeader()
		should.Equal(protocol.TypeI64, elemType)
		should.Equal(3, length)
		should.Equal(uint64(1), iter.ReadUint64())
		should.Equal(uint64(2), iter.ReadUint64())
		should.Equal(uint64(3), iter.ReadUint64())
	}
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

func Test_unmarshal_general_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.I64, 3)
		proto.WriteI64(1)
		proto.WriteI64(2)
		proto.WriteI64(3)
		proto.WriteListEnd()
		var val general.List
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.List{int64(1), int64(2), int64(3)}, val)
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

func Test_marshal_general_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(general.List{
			int64(1), int64(2), int64(3),
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		elemType, length := iter.ReadListHeader()
		should.Equal(protocol.TypeI64, elemType)
		should.Equal(3, length)
		should.Equal(uint64(1), iter.ReadUint64())
		should.Equal(uint64(2), iter.ReadUint64())
		should.Equal(uint64(3), iter.ReadUint64())
	}
}

func Test_marshal_list(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal([]int64{1, 2, 3})
		should.NoError(err)
		iter := c.CreateIterator(output)
		elemType, length := iter.ReadListHeader()
		should.Equal(protocol.TypeI64, elemType)
		should.Equal(3, length)
		should.Equal(uint64(1), iter.ReadUint64())
		should.Equal(uint64(2), iter.ReadUint64())
		should.Equal(uint64(3), iter.ReadUint64())
	}
}
