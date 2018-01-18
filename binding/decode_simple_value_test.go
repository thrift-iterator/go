package binding

import (
	"testing"
	"github.com/v2pro/wombat/generic"
	"github.com/v2pro/wombat"
	"reflect"
	"github.com/thrift-iterator/go/protocol/binary"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/stretchr/testify/require"
)

func init() {
	generic.DynamicCompilationEnabled = true
}

func Test_decode_int64(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteI64(-1)
	iter := binary.NewIterator(buf.Bytes())
	funcObj := generic.Expand(decodeSimpleValue,
		"ST", reflect.TypeOf((*binary.Iterator)(nil)),
			"DT", reflect.PtrTo(wombat.Int64))
	f := funcObj.(func(*error, *int64, *binary.Iterator))
	var err error
	var val int64
	f(&err, &val, iter)
	should.NoError(err)
	should.Equal(int64(-1), val)
}