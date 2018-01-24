package test

//func Test_decode_struct_of_raw_message(t *testing.T) {
//	should := require.New(t)
//	api := thrifter.Config{Protocol: thrifter.ProtocolBinary, DynamicCodegen: true}.Froze()
//	stream := api.NewStream(nil, nil)
//	stream.WriteStruct(general.Struct{
//		0: general.Map{
//			"key1": "value1",
//		},
//		1: "hello",
//	})
//	should.Nil(stream.Error())
//	rawStruct := map[protocol.FieldId]thrifter.RawMessage{}
//	should.NoError(api.Unmarshal(stream.Buffer(), &rawStruct))
//	// parse arg1
//	var arg1 string
//	should.NoError(api.Unmarshal(rawStruct[protocol.FieldId(1)].Buffer, &arg1))
//	should.Equal("hello", arg1)
//	// parse arg0
//	var arg0 map[string]string
//	should.NoError(api.Unmarshal(rawStruct[protocol.FieldId(0)].Buffer, &arg0))
//	should.Equal(map[string]string{"key1": "value1"}, arg0)
//	// modify arg0
//	arg0["key2"] = "value2"
//	encodedArg0, err := api.Marshal(arg0)
//	should.NoError(err)
//	// set arg0 back
//	rawStruct[protocol.FieldId(0)] = thrifter.RawMessage{
//		Buffer: encodedArg0,
//		ThriftType: protocol.TypeMap,
//	}
//	encodedArgs, err := api.Marshal(rawStruct)
//	should.NoError(err)
//	// verify it is changed
//	iter := api.NewIterator(nil, encodedArgs)
//	fmt.Println(iter.ReadStruct())
//}