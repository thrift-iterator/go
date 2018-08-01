package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/general"
)

func Test_skip_map_of_string_key(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.STRING, thrift.I64, 1)
		proto.WriteString("1")
		proto.WriteI64(1)
		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMap(nil))
	}
}

func Test_skip_map_of_string_elem(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.I64, thrift.STRING, 1)
		proto.WriteI64(1)
		proto.WriteString("1")
		proto.WriteMapEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMap(nil))
	}
}

func Test_unmarshal_general_map_of_string_key(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.STRING, thrift.I64, 1)
		proto.WriteString("1")
		proto.WriteI64(1)
		proto.WriteMapEnd()
		var val general.Map
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.Map{
			"1": int64(1),
		}, val)
	}
}

func Test_unmarshal_map_of_string_key(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMapBegin(thrift.STRING, thrift.I64, 1)
		proto.WriteString("1")
		proto.WriteI64(1)
		proto.WriteMapEnd()
		var val map[string]int64
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(map[string]int64{
			"1": 1,
		}, val)
	}
}

func Test_marshal_general_map_of_string_key(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		m := general.Map{
			"1": int64(1),
		}

		output, err := c.Marshal(m)
		should.NoError(err)
		output1, err := c.Marshal(&m)
		should.NoError(err)
		should.Equal(output, output1)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.Map{
			"1": int64(1),
		}, val)
	}
}

func Test_marshal_map_of_string_key(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		m := map[string]int64{
			"1": 1,
		}

		output, err := c.Marshal(m)
		should.NoError(err)
		output1, err := c.Marshal(&m)
		should.NoError(err)
		should.Equal(output, output1)
		var val general.Map
		should.NoError(c.Unmarshal(output, &val))
		should.Equal(general.Map{
			"1": int64(1),
		}, val)
	}
}
