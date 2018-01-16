package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
)

func Test_decode_map(t *testing.T) {
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
	iter := thrifter.NewIterator(buf.Bytes())
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
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipMap())
}