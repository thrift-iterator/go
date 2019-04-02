package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test/api/binding_test"
)

func Test_binding(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	transport := thrift.NewTFramedTransport(buf)
	proto := thrift.NewTBinaryProtocol(transport, true, true)
	proto.WriteStructBegin("hello")
	proto.WriteFieldBegin("field1", thrift.I64, 1)
	proto.WriteI64(1024)
	proto.WriteFieldEnd()
	proto.WriteFieldStop()
	proto.WriteStructEnd()
	transport.Flush()
	var val binding_test.TestObject
	should.NoError(api.Unmarshal(buf.Bytes()[4:], &val))
	should.Equal(int64(1024), val.Field1)
}
