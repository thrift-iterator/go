# thrifter

decode/encode thrift message without IDL

Why?

* IDL generated model is ugly, we do not want *int as my field type. We like how json.Unmarshal works, binds data to existing go struct.
* original thrift library for go is slow, verified in several benchmarks. It is even slower than [json-iterator](https://github.com/json-iterator/go)
* thrift proxy can not load all possible IDL in advance, we need to decode/encode in a generic way to modify embedded header.

# unmarshal

```go
func Test_unmarshal_message(t *testing.T) {
	should := require.New(t)
	input, err := hex.DecodeString("800100010000000568656c6c6f0000000c0b00010000000a73657373696f6e2d69640c00020c00010a000100000000000000010a000200000000000000000b00030000000f43616c6c46726f6d496e626f756e64000c00020b0001000000093132372e302e302e310a000200000000000004d2000b00030000000568656c6c6f000c00030c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000f00040c000000010c00020c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000000")
	should.NoError(err)
	var msg protocol.Message
	err = thrifter.Unmarshal(input, &msg)
	should.NoError(err)
	fmt.Println(msg.MessageType)
	fmt.Println(msg.MessageName)
	for fieldId, fieldValue := range msg.Arguments {
		fmt.Println(fieldId, fieldValue)
	}
}
```

the output is

```
1
hello
1 session-id
2 map[1:map[1:1 2:0 3:CallFromInbound] 2:map[1:127.0.0.1 2:1234] 3:hello]
3 map[1:map[1:2 2:0 3:ReturnInbound] 2:world]
4 [map[2:map[1:map[1:2 2:0 3:ReturnInbound] 2:world]]]
```

the mapping between thrift type and go type

* LIST => `[]interface{}`
* MAP => `map[interface{}]interface{}`
* STRUCT => `map[protocol.FieldId]interface{}`

# marshal

```go
func Test_marshal_message(t *testing.T) {
	should := require.New(t)
	msg := protocol.Message{
		MessageHeader: protocol.MessageHeader{
			Version: protocol.VERSION_1,
			MessageType: protocol.CALL,
			MessageName: "hello",
			SeqId: protocol.SeqId(17),
		},
		Arguments: map[protocol.FieldId]interface{} {
			protocol.FieldId(1): int64(1),
			protocol.FieldId(2): int64(2),
		},
	}
	output, err := thrifter.Marshal(msg)
	should.Nil(err)
	var msgRead protocol.Message
	err = thrifter.Unmarshal(output, &msgRead)
	should.NoError(err)
	fmt.Println(msgRead.MessageType)
	fmt.Println(msgRead.MessageName)
	for fieldId, fieldValue := range msgRead.Arguments {
		fmt.Println(fieldId, fieldValue)
	}
}
```

the output is 

```
1
hello
2 2
1 1
```
