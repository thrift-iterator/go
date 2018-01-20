package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/test"
)

func Test_unmarshal_ptr_int64(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteI64(2)
		proto.WriteListEnd()
		var val *int64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(int64(2), *val)
	}
}