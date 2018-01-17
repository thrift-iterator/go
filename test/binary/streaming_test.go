package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"encoding/hex"
	"bytes"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/protocol"
	"fmt"
)

func Test_decode_framed_transport(t *testing.T) {
	should := require.New(t)
	input, err := hex.DecodeString("800100010000000568656c6c6f0000000c0b00010000000a73657373696f6e2d69640c00020c00010a000100000000000000010a000200000000000000000b00030000000f43616c6c46726f6d496e626f756e64000c00020b0001000000093132372e302e302e310a000200000000000004d2000b00030000000568656c6c6f000c00030c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000f00040c000000010c00020c00010a000100000000000000020a000200000000000000000b00030000000d52657475726e496e626f756e64000b000200000005776f726c64000000")
	should.NoError(err)
	size := len(input)
	input = append([]byte{
		byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
	}, input...)
	reader := bytes.NewBuffer(input)
	decoder := thrifter.NewDecoder(reader)
	var msg protocol.Message
	should.NoError(decoder.Decode(&msg))
	fmt.Println(msg.MessageType)
	fmt.Println(msg.MessageName)
	for fieldId, fieldValue := range msg.Arguments {
		fmt.Println(fieldId, fieldValue)
	}
}
