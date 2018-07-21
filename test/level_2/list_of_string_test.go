package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/general"
)

func Test_skip_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRING, 3)
		proto.WriteString("a")
		proto.WriteString("b")
		proto.WriteString("c")
		proto.WriteListEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipList(nil))
	}
}

func Test_unmarshal_general_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRING, 3)
		proto.WriteString("a")
		proto.WriteString("b")
		proto.WriteString("c")
		proto.WriteListEnd()
		var val general.List
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal(general.List{"a", "b", "c"}, val)
	}
}

func Test_unmarshal_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.UnmarshalCombinations {
		buf, proto := c.CreateProtocol()
		proto.WriteListBegin(thrift.STRING, 3)
		proto.WriteString("a")
		proto.WriteString("b")
		proto.WriteString("c")
		proto.WriteListEnd()
		var val []string
		should.NoError(c.Unmarshal(buf.Bytes(), &val))
		should.Equal([]string{
			"a", "b", "c",
		}, val)
	}
}

func Test_marshal_general_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		lst := general.List{
			"a", "b", "c",
		}

		output, err := c.Marshal(lst)
		should.NoError(err)
		output1, err := c.Marshal(&lst)
		should.NoError(err)
		should.Equal(output, output1)
		var val []string
		should.NoError(c.Unmarshal(output, &val))
		should.Equal([]string{
			"a", "b", "c",
		}, val)
	}
}

func Test_marshal_list_of_string(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations {
		lst := []string{"a", "b", "c"}

		output, err := c.Marshal(lst)
		should.NoError(err)
		output1, err := c.Marshal(&lst)
		should.Equal(output, output1)
		should.NoError(err)
		var val []string
		should.NoError(c.Unmarshal(output, &val))
		should.Equal([]string{
			"a", "b", "c",
		}, val)
	}
}
