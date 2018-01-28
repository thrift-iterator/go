package test

import (
	"testing"
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go"
	"fmt"
	"github.com/thrift-iterator/go/protocol"
	"bytes"
	"github.com/thrift-iterator/go/general"
)

type combination struct {
	encoded string
	api     thrifter.API
}

var combinations = []combination{
	{
		encoded: "800100010000000568656c6c6f0000000c0b00010000000a73657373696f6e2d69640c00020c00010a000100000000000000010a000200000000000000000b00030000000f43616c6c46726f6d496e626f756e64000c00020b0001000000093132372e302e302e310a000200000000000004d2000b00030000000568656c6c6f000c00030c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000f00040c000000010c00020c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000000",
		api:     thrifter.Config{Protocol: thrifter.ProtocolBinary}.Froze(),
	},
	{
		encoded: "82210c0568656c6c6f180a73657373696f6e2d69641c1c16021600180f43616c6c46726f6d496e626f756e64001c18093132372e302e302e3116a41300180568656c6c6f001c1c16041600180d52657475726e496e626f756e64001805776f726c6400191c2c1c16041600180d52657475726e496e626f756e64001805776f726c64000000",
		api:     thrifter.Config{Protocol: thrifter.ProtocolCompact}.Froze(),
	},
}

func Test_unmarshal_message(t *testing.T) {
	should := require.New(t)
	for _, c := range combinations {
		input, err := hex.DecodeString(c.encoded)
		should.NoError(err)
		var msg general.Message
		err = c.api.Unmarshal(input, &msg)
		should.NoError(err)
		fmt.Println(msg.MessageType)
		fmt.Println(msg.MessageName)
		for fieldId, fieldValue := range msg.Arguments {
			fmt.Println("!!!", fieldId, fieldValue)
		}
	}
}

func Test_marshal_message(t *testing.T) {
	should := require.New(t)
	msg := general.Message{
		MessageHeader: protocol.MessageHeader{
			MessageType: protocol.MessageTypeCall,
			MessageName: "hello",
			SeqId:       protocol.SeqId(17),
		},
		Arguments: general.Struct{
			protocol.FieldId(1): int64(1),
			protocol.FieldId(2): int64(2),
		},
	}
	output, err := thrifter.Marshal(msg)
	should.Nil(err)
	var msgRead general.Message
	err = thrifter.Unmarshal(output, &msgRead)
	should.NoError(err)
	fmt.Println(msgRead.MessageType)
	fmt.Println(msgRead.MessageName)
	for fieldId, fieldValue := range msgRead.Arguments {
		fmt.Println(fieldId, fieldValue)
	}
}

func Test_decode_message(t *testing.T) {
	should := require.New(t)
	input, err := hex.DecodeString("800100010000000568656c6c6f0000000c0b00010000000a73657373696f6e2d69640c00020c00010a000100000000000000010a000200000000000000000b00030000000f43616c6c46726f6d496e626f756e64000c00020b0001000000093132372e302e302e310a000200000000000004d2000b00030000000568656c6c6f000c00030c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000f00040c000000010c00020c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000000")
	should.NoError(err)
	reader := bytes.NewBuffer(input)
	cfg := thrifter.Config{Protocol: thrifter.ProtocolBinary}.Froze()
	decoder := cfg.NewDecoder(reader, nil)
	var msg general.Message
	should.NoError(decoder.Decode(&msg))
	fmt.Println(msg.MessageType)
	fmt.Println(msg.MessageName)
	for fieldId, fieldValue := range msg.Arguments {
		fmt.Println(fieldId, fieldValue)
	}
}

func Test_encode_message(t *testing.T) {
	should := require.New(t)
	msg := general.Message{
		MessageHeader: protocol.MessageHeader{
			MessageType: protocol.MessageTypeCall,
			MessageName: "hello",
			SeqId:       protocol.SeqId(17),
		},
		Arguments: general.Struct{
			protocol.FieldId(1): int64(1),
			protocol.FieldId(2): int64(2),
		},
	}
	var msgRead general.Message
	buf := bytes.NewBuffer(nil)
	cfg := thrifter.Config{Protocol: thrifter.ProtocolBinary}.Froze()
	encoder := cfg.NewEncoder(buf)
	should.NoError(encoder.Encode(msg))
	err := cfg.Unmarshal(buf.Bytes(), &msgRead)
	should.NoError(err)
	fmt.Println(msgRead.MessageType)
	fmt.Println(msgRead.MessageName)
	for fieldId, fieldValue := range msgRead.Arguments {
		fmt.Println(fieldId, fieldValue)
	}
}
