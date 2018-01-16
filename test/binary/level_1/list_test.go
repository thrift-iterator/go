package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
)

func Test_decode_list(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.I64, 3)
	proto.WriteI64(1)
	proto.WriteI64(2)
	proto.WriteI64(3)
	proto.WriteListEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	elemType, length := iter.ReadList()
	should.Equal(protocol.I64, elemType)
	should.Equal(3, length)
	should.Equal(uint64(1), iter.ReadUInt64())
	should.Equal(uint64(2), iter.ReadUInt64())
	should.Equal(uint64(3), iter.ReadUInt64())
}

func Test_skip_list(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteListBegin(thrift.I64, 3)
	proto.WriteI64(1)
	proto.WriteI64(2)
	proto.WriteI64(3)
	proto.WriteListEnd()
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipList())
}
