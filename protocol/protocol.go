package protocol

// Type constants in the Thrift protocol
type TType byte
type TMessageType int32
type SeqId int32
type FieldId int16
type Version uint32

const VERSION_1 Version = 0x80010000

const (
	INVALID_TMESSAGE_TYPE TMessageType = 0
	CALL                  TMessageType = 1
	REPLY                 TMessageType = 2
	EXCEPTION             TMessageType = 3
	ONEWAY                TMessageType = 4
)

const (
	STOP   TType = 0
	VOID   TType = 1
	BOOL   TType = 2
	BYTE   TType = 3
	I08    TType = 3
	DOUBLE TType = 4
	I16    TType = 6
	I32    TType = 8
	I64    TType = 10
	STRING TType = 11
	UTF7   TType = 11
	STRUCT TType = 12
	MAP    TType = 13
	SET    TType = 14
	LIST   TType = 15
	UTF8   TType = 16
	UTF16  TType = 17
	//BINARY = 18   wrong and unusued
)

var typeNames = map[TType]string{
	STOP:   "STOP",
	VOID:   "VOID",
	BOOL:   "BOOL",
	BYTE:   "BYTE",
	DOUBLE: "DOUBLE",
	I16:    "I16",
	I32:    "I32",
	I64:    "I64",
	STRING: "STRING",
	STRUCT: "STRUCT",
	MAP:    "MAP",
	SET:    "SET",
	LIST:   "LIST",
	UTF8:   "UTF8",
	UTF16:  "UTF16",
}

func (p TType) String() string {
	if s, ok := typeNames[p]; ok {
		return s
	}
	return "Unknown"
}

type MessageHeader struct {
	Version     Version
	MessageName string
	MessageType TMessageType
	SeqId       SeqId
}

type Message struct {
	MessageHeader
	Arguments map[FieldId]interface{}
}

func (msg *Message) Set(msgRead *Message) {
	msg.Version = msgRead.Version
	msg.MessageType = msgRead.MessageType
	msg.MessageName = msgRead.MessageName
	msg.SeqId = msgRead.SeqId
	msg.Arguments = msgRead.Arguments
}
