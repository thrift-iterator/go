package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/general"
	"github.com/thrift-iterator/go/raw"
)

func Test_decode_map_by_iterator(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.STRING, thrift.I64, 3)
		proto.WriteString("k1")
		proto.WriteI64(1)
		proto.WriteString("k2")
		proto.WriteI64(2)
		proto.WriteString("k3")
		proto.WriteI64(3)
		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		keyType, elemType, length := iter.ReadMapHeader()
		should.Equal(protocol.TypeString, keyType)
		should.Equal(protocol.TypeI64, elemType)
		should.Equal(3, length)
		should.Equal("k1", iter.ReadString())
		should.Equal(uint64(1), iter.ReadUint64())
		should.Equal("k2", iter.ReadString())
		should.Equal(uint64(2), iter.ReadUint64())
		should.Equal("k3", iter.ReadString())
		should.Equal(uint64(3), iter.ReadUint64())
	}
}

func Test_encode_map_by_stream(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		stream := c.CreateStream()
		stream.WriteMapHeader(protocol.TypeString, protocol.TypeI64, 3)
		stream.WriteString("k1")
		stream.WriteUint64(1)
		stream.WriteString("k2")
		stream.WriteUint64(2)
		stream.WriteString("k3")
		stream.WriteUint64(3)
		iter := c.CreateIterator(stream.Buffer())
		keyType, elemType, length := iter.ReadMapHeader()
		should.Equal(protocol.TypeString, keyType)
		should.Equal(protocol.TypeI64, elemType)
		should.Equal(3, length)
		should.Equal("k1", iter.ReadString())
		should.Equal(uint64(1), iter.ReadUint64())
		should.Equal("k2", iter.ReadString())
		should.Equal(uint64(2), iter.ReadUint64())
		should.Equal("k3", iter.ReadString())
		should.Equal(uint64(3), iter.ReadUint64())
	}
}

func Test_skip_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 3)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteI32(3)
		proto.WriteI64(3)
		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMap(nil))
	}
}

func Test_unmarshal_general_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 3)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteI32(3)
		proto.WriteI64(3)
		proto.WriteMapEnd()
		var val general.Map
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.Map{
			int32(1): int64(1),
			int32(2): int64(2),
			int32(3): int64(3),
		}, val)
	}
}


func Test_unmarshal_raw_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 3)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteI32(3)
		proto.WriteI64(3)
		proto.WriteMapEnd()
		var val raw.Map
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(3, len(val.Entries))
		should.Equal(protocol.TypeI32, val.KeyType)
		should.Equal(protocol.TypeI64, val.ElementType)
		iter := c.CreateIterator(val.Entries[int32(1)].Element)
		should.Equal(int64(1), iter.ReadInt64())
	}
}

func Test_unmarshal_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 3)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteI32(3)
		proto.WriteI64(3)
		proto.WriteMapEnd()
		val := map[int32]int64{}
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(map[int32]int64{
			int32(1): int64(1),
			int32(2): int64(2),
			int32(3): int64(3),
		}, val)
	}
}

func Test_marshal_general_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(general.Map{
			"k1": int64(1),
			"k2": int64(2),
			"k3": int64(3),
		})
		should.NoError(err)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.Map{
			int32(1): int64(1),
			int32(2): int64(2),
			int32(3): int64(3),
		}, val)
	}
}


func Test_marshal_raw_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I32, thrift.I64, 3)
		proto.WriteI32(1)
		proto.WriteI64(1)
		proto.WriteI32(2)
		proto.WriteI64(2)
		proto.WriteI32(3)
		proto.WriteI64(3)
		proto.WriteMapEnd()
		var val raw.Map
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		output, err := c.Marshal(val)
		should.NoError(err)
		var generalVal general.Map
		should.NoError(c.Unmarshal(output, &generalVal))
		should.Equal(general.Map{
			int32(1): int64(1),
			int32(2): int64(2),
			int32(3): int64(3),
		}, generalVal)
	}
}

func Test_marshal_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		output, err := c.Marshal(map[string]int64{
			"k1": int64(1),
			"k2": int64(2),
			"k3": int64(3),
		})
		should.NoError(err)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.Map{
			int32(1): int64(1),
			int32(2): int64(2),
			int32(3): int64(3),
		}, val)
	}
}
