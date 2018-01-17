package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
)

func Test_decode_map_by_iterator(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteMapBegin(thrift.STRING, thrift.I64, 3)
	proto.WriteString("k1")
	proto.WriteI64(1)
	proto.WriteString("k2")
	proto.WriteI64(2)
	proto.WriteString("k3")
	proto.WriteI64(3)
	proto.WriteMapEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	keyType, elemType, length := iter.ReadMapHeader()
	should.Equal(protocol.STRING, keyType)
	should.Equal(protocol.I64, elemType)
	should.Equal(3, length)
	should.Equal("k1", iter.ReadString())
	should.Equal(uint64(1), iter.ReadUInt64())
	should.Equal("k2", iter.ReadString())
	should.Equal(uint64(2), iter.ReadUInt64())
	should.Equal("k3", iter.ReadString())
	should.Equal(uint64(3), iter.ReadUInt64())
}

func Test_encode_map_by_stream(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteMapHeader(protocol.STRING, protocol.I64, 3)
	stream.WriteString("k1")
	stream.WriteUInt64(1)
	stream.WriteString("k2")
	stream.WriteUInt64(2)
	stream.WriteString("k3")
	stream.WriteUInt64(3)
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	keyType, elemType, length := iter.ReadMapHeader()
	should.Equal(protocol.STRING, keyType)
	should.Equal(protocol.I64, elemType)
	should.Equal(3, length)
	should.Equal("k1", iter.ReadString())
	should.Equal(uint64(1), iter.ReadUInt64())
	should.Equal("k2", iter.ReadString())
	should.Equal(uint64(2), iter.ReadUInt64())
	should.Equal("k3", iter.ReadString())
	should.Equal(uint64(3), iter.ReadUInt64())
}

func Test_decode_map_as_object(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteMapBegin(thrift.STRING, thrift.I64, 3)
	proto.WriteString("k1")
	proto.WriteI64(1)
	proto.WriteString("k2")
	proto.WriteI64(2)
	proto.WriteString("k3")
	proto.WriteI64(3)
	proto.WriteMapEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	obj := iter.ReadMap()
	should.Equal(map[interface{}]interface{}{
		"k1": int64(1),
		"k2": int64(2),
		"k3": int64(3),
	}, obj)
}

func Test_encode_map_from_object(t *testing.T) {
	should := require.New(t)
	stream := thrifter.NewBufferedStream(nil)
	stream.WriteMap(map[interface{}]interface{}{
		"k1": int64(1),
		"k2": int64(2),
		"k3": int64(3),
	})
	iter := thrifter.NewBufferedIterator(stream.Buffer())
	obj := iter.ReadMap()
	should.Equal(map[interface{}]interface{}{
		"k1": int64(1),
		"k2": int64(2),
		"k3": int64(3),
	}, obj)
}

func Test_skip_map(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteMapBegin(thrift.I32, thrift.I64, 3)
	proto.WriteI32(1)
	proto.WriteI64(1)
	proto.WriteI32(2)
	proto.WriteI64(2)
	proto.WriteI32(3)
	proto.WriteI64(3)
	proto.WriteMapEnd()
	iter := thrifter.NewBufferedIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipMap(nil))
}
