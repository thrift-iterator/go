package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/general"
)

func Test_skip_map_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.MAP, 1)
		proto.WriteI64(1)

		proto.WriteMapBegin(thrift.STRING, thrift.I64, 1)
		proto.WriteString("k1")
		proto.WriteI64(1)
		proto.WriteMapEnd()

		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMap(nil))
	}
}

func Test_unmarshal_general_map_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.MAP, 1)
		proto.WriteI64(1)

		proto.WriteMapBegin(thrift.STRING, thrift.I64, 1)
		proto.WriteString("k1")
		proto.WriteI64(1)
		proto.WriteMapEnd()

		proto.WriteMapEnd()
		var val general.Map
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.Map{
			"k1": int64(1),
		}, val[int64(1)])
	}
}

func Test_marshal_general_map_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		m := general.Map{
			int64(1): general.Map{
				"k1": int64(1),
			},
		}

		output, err := c.Marshal(m)
		should.NoError(err)
		output1, err := c.Marshal(&m)
		should.NoError(err)
		should.Equal(output, output1)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.Map{
			"k1": int64(1),
		}, val[int64(1)])
	}
}

func Test_marshal_map_of_map(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		m := map[int64]map[string]int64{
			1: {"k1": 1},
		}

		output, err := c.Marshal(m)
		should.NoError(err)
		output1, err := c.Marshal(&m)
		should.NoError(err)
		should.Equal(output, output1)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.Map{
			"k1": int64(1),
		}, val[int64(1)])
	}
}