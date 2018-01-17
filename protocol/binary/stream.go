package binary

import (
	"math"
	"github.com/thrift-iterator/go/protocol"
	"fmt"
	"io"
)

type Stream struct {
	writer io.Writer
	buf    []byte
	err    error
}

func NewStream(writer io.Writer, buf []byte) *Stream {
	return &Stream{
		writer: writer,
		buf:    buf,
	}
}

func (stream *Stream) Error() error {
	return stream.err
}

func (stream *Stream) ReportError(operation string, err string) {
	if stream.err == nil {
		stream.err = fmt.Errorf("%s: %s", operation, err)
	}
}

func (stream *Stream) Buffer() []byte {
	return stream.buf
}

func (stream *Stream) Reset(writer io.Writer) {
	stream.writer = writer
	stream.err = nil
	stream.buf = stream.buf[:0]
}

func (stream *Stream) Flush() {
	if stream.writer == nil {
		return
	}
	_, err := stream.writer.Write(stream.buf)
	if err != nil {
		stream.ReportError("Flush", err.Error())
		return
	}
	stream.buf = stream.buf[:0]
}

func (stream *Stream) WriteMessageHeader(header protocol.MessageHeader) {
	versionAndMessageType := uint32(header.Version) | uint32(header.MessageType)
	stream.WriteUInt32(versionAndMessageType)
	stream.WriteString(header.MessageName)
	stream.WriteInt32(int32(header.SeqId))
}

func (stream *Stream) WriteMessage(message protocol.Message) {
	stream.WriteMessageHeader(message.MessageHeader)
	stream.WriteStruct(message.Arguments)
}

func (stream *Stream) WriteListHeader(elemType protocol.TType, length int) {
	stream.buf = append(stream.buf, byte(elemType),
		byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
}

func (stream *Stream) WriteList(val []interface{}) {
	if len(val) == 0 {
		stream.ReportError("WriteList", "input is empty slice, can not tell element type")
		return
	}
	elemType, elemWriter := stream.WriterOf(val[0])
	stream.WriteListHeader(elemType, len(val))
	for _, elem := range val {
		elemWriter(elem)
	}
}

func (stream *Stream) WriteStructField(fieldType protocol.TType, fieldId protocol.FieldId) {
	stream.buf = append(stream.buf, byte(fieldType), byte(fieldId>>8), byte(fieldId))
}

func (stream *Stream) WriteStructFieldStop() {
	stream.buf = append(stream.buf, byte(protocol.STOP))
}

func (stream *Stream) WriteStruct(val map[protocol.FieldId]interface{}) {
	for key, elem := range val {
		switch typedElem := elem.(type) {
		case bool:
			stream.WriteStructField(protocol.BOOL, key)
			stream.WriteBool(typedElem)
		case int8:
			stream.WriteStructField(protocol.I08, key)
			stream.WriteInt8(typedElem)
		case uint8:
			stream.WriteStructField(protocol.I08, key)
			stream.WriteUInt8(typedElem)
		case int16:
			stream.WriteStructField(protocol.I16, key)
			stream.WriteInt16(typedElem)
		case uint16:
			stream.WriteStructField(protocol.I16, key)
			stream.WriteUInt16(typedElem)
		case int32:
			stream.WriteStructField(protocol.I32, key)
			stream.WriteInt32(typedElem)
		case uint32:
			stream.WriteStructField(protocol.I32, key)
			stream.WriteUInt32(typedElem)
		case int64:
			stream.WriteStructField(protocol.I64, key)
			stream.WriteInt64(typedElem)
		case uint64:
			stream.WriteStructField(protocol.I64, key)
			stream.WriteUInt64(typedElem)
		case float64:
			stream.WriteStructField(protocol.DOUBLE, key)
			stream.WriteFloat64(typedElem)
		case string:
			stream.WriteStructField(protocol.STRING, key)
			stream.WriteString(typedElem)
		case []interface{}:
			stream.WriteStructField(protocol.LIST, key)
			stream.WriteList(typedElem)
		case map[interface{}]interface{}:
			stream.WriteStructField(protocol.MAP, key)
			stream.WriteMap(typedElem)
		case map[protocol.FieldId]interface{}:
			stream.WriteStructField(protocol.STRUCT, key)
			stream.WriteStruct(typedElem)
		default:
			panic("unsupported type")
		}
	}
	stream.WriteStructFieldStop()
	stream.Flush()
}

func (stream *Stream) WriteMapHeader(keyType protocol.TType, elemType protocol.TType, length int) {
	stream.buf = append(stream.buf, byte(keyType), byte(elemType),
		byte(length>>24), byte(length>>16), byte(length>>8), byte(length))
}

func (stream *Stream) WriteMap(val map[interface{}]interface{}) {
	hasSample, sampleKey, sampleElem := takeSampleFromMap(val)
	if !hasSample {
		stream.ReportError("WriteMap", "input is empty map, can not tell element type")
		return
	}
	keyType, keyWriter := stream.WriterOf(sampleKey)
	elemType, elemWriter := stream.WriterOf(sampleElem)
	stream.WriteMapHeader(keyType, elemType, len(val))
	for key, elem := range val {
		keyWriter(key)
		elemWriter(elem)
	}
}

func takeSampleFromMap(val map[interface{}]interface{}) (bool, interface{}, interface{}) {
	for key, elem := range val {
		return true, key, elem
	}
	return false, nil, nil
}

func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.WriteUInt8(1)
	} else {
		stream.WriteUInt8(0)
	}
}

func (stream *Stream) WriteInt8(val int8) {
	stream.WriteUInt8(uint8(val))
}

func (stream *Stream) WriteUInt8(val uint8) {
	stream.buf = append(stream.buf, byte(val))
}

func (stream *Stream) WriteInt16(val int16) {
	stream.WriteUInt16(uint16(val))
}

func (stream *Stream) WriteUInt16(val uint16) {
	stream.buf = append(stream.buf, byte(val>>8), byte(val))
}

func (stream *Stream) WriteInt32(val int32) {
	stream.WriteUInt32(uint32(val))
}

func (stream *Stream) WriteUInt32(val uint32) {
	stream.buf = append(stream.buf, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

func (stream *Stream) WriteInt64(val int64) {
	stream.WriteUInt64(uint64(val))
}

func (stream *Stream) WriteUInt64(val uint64) {
	stream.buf = append(stream.buf,
		byte(val>>56), byte(val>>48), byte(val>>40), byte(val>>32),
		byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

func (stream *Stream) WriteFloat64(val float64) {
	stream.WriteUInt64(math.Float64bits(val))
}

func (stream *Stream) WriteBinary(val []byte) {
	stream.WriteUInt32(uint32(len(val)))
	stream.buf = append(stream.buf, val...)
}

func (stream *Stream) WriteString(val string) {
	stream.WriteUInt32(uint32(len(val)))
	stream.buf = append(stream.buf, val...)
}

func (stream *Stream) WriterOf(sample interface{}) (protocol.TType, func(interface{})) {
	switch sample.(type) {
	case bool:
		return protocol.BOOL, func(val interface{}) {
			stream.WriteBool(val.(bool))
		}
	case int8:
		return protocol.I08, func(val interface{}) {
			stream.WriteInt8(val.(int8))
		}
	case uint8:
		return protocol.I08, func(val interface{}) {
			stream.WriteUInt8(val.(uint8))
		}
	case int16:
		return protocol.I16, func(val interface{}) {
			stream.WriteInt16(val.(int16))
		}
	case uint16:
		return protocol.I16, func(val interface{}) {
			stream.WriteUInt16(val.(uint16))
		}
	case int32:
		return protocol.I32, func(val interface{}) {
			stream.WriteInt32(val.(int32))
		}
	case uint32:
		return protocol.I32, func(val interface{}) {
			stream.WriteUInt32(val.(uint32))
		}
	case int64:
		return protocol.I64, func(val interface{}) {
			stream.WriteInt64(val.(int64))
		}
	case uint64:
		return protocol.I64, func(val interface{}) {
			stream.WriteUInt64(val.(uint64))
		}
	case float64:
		return protocol.DOUBLE, func(val interface{}) {
			stream.WriteFloat64(val.(float64))
		}
	case string:
		return protocol.STRING, func(val interface{}) {
			stream.WriteString(val.(string))
		}
	case []interface{}:
		return protocol.LIST, func(val interface{}) {
			stream.WriteList(val.([]interface{}))
		}
	case map[interface{}]interface{}:
		return protocol.MAP, func(val interface{}) {
			stream.WriteMap(val.(map[interface{}]interface{}))
		}
	case map[protocol.FieldId]interface{}:
		return protocol.STRUCT, func(val interface{}) {
			stream.WriteStruct(val.(map[protocol.FieldId]interface{}))
		}
	default:
		panic("unsupported type")
	}
}
