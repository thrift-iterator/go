package binding

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/v2pro/wombat/generic"
	"reflect"
	"github.com/v2pro/wombat"
	"github.com/thrift-iterator/go/protocol/binary"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func Test_decode_anything(t *testing.T) {
	should := require.New(t)
	buf := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocol(buf, true, true)
	proto.WriteI64(-1)
	iter := binary.NewIterator(buf.Bytes())
	funcObj := generic.Expand(DecodeAnything,
		"ST", reflect.TypeOf((*binary.Iterator)(nil)),
		"DT", reflect.PtrTo(wombat.Int64))
	f := funcObj.(func(interface{}, interface{}))
	var val int64
	f(&val, iter)
	should.Equal(int64(-1), val)
}
