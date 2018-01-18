package test

import (
	"testing"
	"github.com/v2pro/wombat/generic"
	"github.com/v2pro/wombat"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go"
)

func init() {
	generic.DynamicCompilationEnabled = true
}

func Test_binding(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteI64(-1)
	var val int64
	api := thrifter.Config{
		Protocol: thrifter.ProtocolBinary,
	}.Decode(wombat.Int64).Froze()
	should.NoError(api.Unmarshal(buf.Bytes(), &val))
	should.Equal(int64(-1), val)
}
