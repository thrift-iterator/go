package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_skip_list_of_map(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
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
	iter := thrifter.NewIterator(buf.Bytes())
	should.Equal(buf.Bytes(), iter.SkipList())
}