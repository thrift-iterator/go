# thrifter

decode/encode thrift message without IDL

Why?

* because IDL generated model is ugly and inflexible, it is seldom used in application directly. instead we define another model, which leads to bad performance.
  * bytes need to be copied twice 
  * more objects to gc
* thrift proxy can not know all possible IDL in advance, we need to decode/encode in a generic way to modify embedded header.
* official thrift library for go is slow, verified in several benchmarks. It is even slower than [json-iterator](https://github.com/json-iterator/go)


# marshal without idl

```go
type SampleArgs struct {
	SessionId string `thrift:",1"`
}

func Test_encode_framed_message_header_and_args(t *testing.T) {
	should := require.New(t)
	msgHeader := protocol.MessageHeader{
		MessageType: protocol.MessageTypeCall,
		MessageName: "hello",
		SeqId:       protocol.SeqId(17),
	}
	var msgRead protocol.Message
	buf := bytes.NewBuffer(nil)
	encoder := thrifter.NewEncoder(buf)
	// write message header
	should.NoError(encoder.Encode(msgHeader))
	// write message args
	should.NoError(encoder.Encode(SampleArgs{
		SessionId: "session-id",
	}))
	err := thrifter.Unmarshal(buf.Bytes(), &msgRead)
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
1 session-id
```

# unmarshal without idl

```go
type SampleArgs struct {
	SessionId string `thrift:",1"`
}

func Test_decode_framed_message_header_and_args(t *testing.T) {
	should := require.New(t)
	input, err := hex.DecodeString("800100010000000568656c6c6f0000000c0b00010000000a73657373696f6e2d69640c00020c00010a000100000000000000010a000200000000000000000b00030000000f43616c6c46726f6d496e626f756e64000c00020b0001000000093132372e302e302e310a000200000000000004d2000b00030000000568656c6c6f000c00030c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000f00040c000000010c00020c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000000")
	should.NoError(err)
	size := len(input)
	input = append([]byte{
		byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
	}, input...)
	reader := bytes.NewBuffer(input)
	decoder := thrifter.NewDecoder(reader, nil)
	// strip out the message header
	var msg protocol.MessageHeader
	should.NoError(decoder.Decode(&msg))
	fmt.Println(msg.MessageType)
	fmt.Println(msg.MessageName)
	// bind args into struct
	var args SampleArgs
	should.NoError(decoder.Decode(&args))
	fmt.Println(args.SessionId)
}
```

the output is

```
1
hello
session-id
```

if not bind to struct, the default mapping between thrift type and go type

* LIST => `[]interface{}`
* MAP => `map[interface{}]interface{}`
* STRUCT => `map[protocol.FieldId]interface{}`
