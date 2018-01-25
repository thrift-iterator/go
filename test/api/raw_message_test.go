package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/general"
	"fmt"
	"github.com/thrift-iterator/go/raw"
	"github.com/thrift-iterator/go/protocol"
)

func Test_decode_struct_of_raw_message(t *testing.T) {
	should := require.New(t)
	api := thrifter.Config{Protocol: thrifter.ProtocolBinary, DynamicCodegen: true}.Froze()
	output, err := api.Marshal(general.Struct{
		0: general.Map{
			"key1": "value1",
		},
		1: "hello",
	})
	should.Nil(err)
	rawStruct := raw.Struct{}
	should.NoError(api.Unmarshal(output, &rawStruct))
	// parse arg1
	var arg1 string
	should.NoError(api.Unmarshal(rawStruct[protocol.FieldId(1)].Buffer, &arg1))
	should.Equal("hello", arg1)
	// parse arg0
	var arg0 map[string]string
	should.NoError(api.Unmarshal(rawStruct[protocol.FieldId(0)].Buffer, &arg0))
	should.Equal(map[string]string{"key1": "value1"}, arg0)
	// modify arg0
	arg0["key2"] = "value2"
	encodedArg0, err := api.Marshal(arg0)
	should.NoError(err)
	// set arg0 back
	rawStruct[protocol.FieldId(0)] = raw.StructField{
		Buffer: encodedArg0,
		Type: protocol.TypeMap,
	}
	encodedArgs, err := api.Marshal(rawStruct)
	should.NoError(err)
	// verify it is changed
	var val general.Struct
	should.NoError(api.Unmarshal(encodedArgs, &val))
	fmt.Println(val)
}