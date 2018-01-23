package test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/test"
)

func Test_skip_message(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMessageBegin("hello", thrift.CALL, 17)
		proto.WriteStructBegin("args")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteFieldBegin("field2", thrift.I64, 2)
		proto.WriteI64(2)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		proto.WriteMessageEnd()
		iter := c.CreateIterator(buf.Bytes())
		should.Equal(buf.Bytes(), iter.SkipMessage(nil))
	}
}

func Test_unmarshal_message(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		buf, proto := c.CreateProtocol()
		proto.WriteMessageBegin("hello", thrift.CALL, 17)
		proto.WriteStructBegin("args")
		proto.WriteFieldBegin("field1", thrift.I64, 1)
		proto.WriteI64(1)
		proto.WriteFieldBegin("field2", thrift.I64, 2)
		proto.WriteI64(2)
		proto.WriteFieldEnd()
		proto.WriteFieldStop()
		proto.WriteStructEnd()
		proto.WriteMessageEnd()
		var msg protocol.Message
		should.NoError(c.Unmarshal(buf.Bytes(), &msg))
		should.Equal("hello", msg.MessageName)
		should.Equal(protocol.MessageTypeCall, msg.MessageType)
		should.Equal(protocol.SeqId(17), msg.SeqId)
		should.Equal(int64(1), msg.Arguments[protocol.FieldId(1)])
		should.Equal(int64(2), msg.Arguments[protocol.FieldId(2)])
	}
}

func Test_marshal_message(t *testing.T) {
	should := require.New(t)
	for _, c := range test.Combinations {
		output, err := c.Marshal(protocol.Message{
			MessageHeader: protocol.MessageHeader{
				MessageType: protocol.MessageTypeCall,
				MessageName: "hello",
				SeqId:       protocol.SeqId(17),
			},
			Arguments: map[protocol.FieldId]interface{}{
				protocol.FieldId(1): int64(1),
				protocol.FieldId(2): int64(2),
			},
		})
		should.NoError(err)
		iter := c.CreateIterator(output)
		msg := iter.ReadMessage()
		should.Equal("hello", msg.MessageName)
		should.Equal(protocol.MessageTypeCall, msg.MessageType)
		should.Equal(protocol.SeqId(17), msg.SeqId)
		should.Equal(int64(1), msg.Arguments[protocol.FieldId(1)])
		should.Equal(int64(2), msg.Arguments[protocol.FieldId(2)])
	}
}
