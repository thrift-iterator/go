package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
	"github.com/thrift-iterator/go/protocol/binary"
	"fmt"
)

type Iterator struct {
	real   *binary.Iterator
	reader io.Reader
	tmp    []byte
	space  []byte
	err    error
}

func NewIterator(reader io.Reader, buf []byte) *Iterator {
	return &Iterator{
		reader: reader, real: binary.NewIterator(nil),
		tmp:    make([]byte, 10),
		space:  buf,
	}
}

func (iter *Iterator) Reset(reader io.Reader, buf []byte) {
	iter.reader = reader
	iter.err = nil
}

func (iter *Iterator) allocate(nBytes int) []byte {
	if len(iter.tmp) < nBytes {
		iter.tmp = make([]byte, nBytes)
	}
	return iter.tmp[:nBytes]
}

func (iter *Iterator) skip(space []byte, nBytes int) []byte {
	tmp := iter.tmp[:nBytes]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("skip", err.Error())
		return nil
	}
	return append(space, tmp...)
}

func (iter *Iterator) Error() error {
	return iter.err
}

func (iter *Iterator) ReportError(operation string, err string) {
	if iter.err == nil {
		iter.err = fmt.Errorf("%s: %s", operation, err)
	}
}

func (iter *Iterator) ReadMessageHeader() protocol.MessageHeader {
	versionAndMessageType := iter.ReadInt32()
	messageType := protocol.TMessageType(versionAndMessageType & 0x0ff)
	version := protocol.Version(int64(int64(versionAndMessageType) & 0xffff0000))
	messageName := iter.ReadString()
	seqId := protocol.SeqId(iter.ReadInt32())
	return protocol.MessageHeader{
		Version:     version,
		MessageName: messageName,
		MessageType: messageType,
		SeqId:       seqId,
	}
}

func (iter *Iterator) ReadMessage() protocol.Message {
	buf := iter.SkipMessage(iter.space[:0])
	if iter.err != nil {
		return protocol.Message{}
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadMessage()
}

func (iter *Iterator) SkipMessage(space []byte) []byte {
	space = iter.skip(space, 4)
	space = iter.SkipBinary(space)
	space = iter.skip(space, 4)
	space = iter.SkipStruct(space)
	return space
}

func (iter *Iterator) ReadStructCB(cb func(fieldType protocol.TType, fieldId protocol.FieldId)) {
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.STOP {
			return
		}
		cb(fieldType, fieldId)
	}
}

func (iter *Iterator) ReadStructHeader() {
	// noop
}

func (iter *Iterator) ReadStructField() (fieldType protocol.TType, fieldId protocol.FieldId) {
	tmp := iter.tmp[:3]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadStructField", err.Error())
		return protocol.STOP, 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadStructField()
}

func (iter *Iterator) ReadStruct() map[protocol.FieldId]interface{} {
	buf := iter.SkipStruct(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadStruct()
}

func (iter *Iterator) SkipStruct(space []byte) []byte {
	for {
		tmp := iter.tmp[:1]
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("SkipStruct", err.Error())
			return nil
		}
		fieldType := protocol.TType(tmp[0])
		space = append(space, tmp[0])
		switch fieldType {
		case protocol.STOP:
			return space
		case protocol.I64, protocol.DOUBLE:
			tmp := iter.tmp[:10]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
		case protocol.LIST:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipList(space)
		case protocol.MAP:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipMap(space)
		case protocol.STRING:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipBinary(space)
		case protocol.STRUCT:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipStruct(space)
		default:
			panic("unsupported type")
		}
	}
}

func (iter *Iterator) ReadListHeader() (elemType protocol.TType, length int) {
	tmp := iter.tmp[:5]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadListHeader", err.Error())
		return protocol.STOP, 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadListHeader()
}

func (iter *Iterator) ReadList() []interface{} {
	buf := iter.SkipList(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadList()
}

func (iter *Iterator) SkipList(space []byte) []byte {
	tmp := iter.tmp[:5]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("SkipList", err.Error())
		return nil
	}
	space = append(space, tmp...)
	iter.real.Reset(nil, tmp)
	elemType, length := iter.real.ReadListHeader()
	switch elemType {
	case protocol.STOP:
		return nil
	case protocol.I64, protocol.DOUBLE:
		tmp := iter.allocate(length * 8)
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("SkipList", err.Error())
			return nil
		}
		space = append(space, tmp...)
		return space
	case protocol.STRING:
		for i := 0; i < length; i++ {
			space = iter.SkipBinary(space)
		}
		return space
	case protocol.LIST:
		for i := 0; i < length; i++ {
			space = iter.SkipList(space)
		}
		return space
	case protocol.MAP:
		for i := 0; i < length; i++ {
			space = iter.SkipMap(space)
		}
		return space
	case protocol.STRUCT:
		for i := 0; i < length; i++ {
			space = iter.SkipStruct(space)
		}
		return space
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) ReadMapHeader() (keyType protocol.TType, elemType protocol.TType, length int) {
	tmp := iter.tmp[:6]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadMapHeader", err.Error())
		return protocol.STOP, protocol.STOP, 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadMapHeader()
}

func (iter *Iterator) ReadMap() map[interface{}]interface{} {
	buf := iter.SkipMap(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(nil, buf)
	return iter.real.ReadMap()
}

func (iter *Iterator) SkipMap(space []byte) []byte {
	tmp := iter.tmp[:6]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("SkipMap", err.Error())
		return nil
	}
	space = append(space, tmp...)
	iter.real.Reset(nil, tmp)
	keyType, elemType, length := iter.real.ReadMapHeader()
	keySize := getTypeSize(keyType)
	elemSize := getTypeSize(elemType)
	if keySize != 0 && elemSize != 0 {
		tmp := iter.allocate(length * (keySize + elemSize))
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("SkipMap", err.Error())
			return nil
		}
		space = append(space, tmp...)
		return space
	}
	var skipKey func(space []byte) []byte
	var skipElem func(space []byte) []byte
	if keySize != 0 {
		skipKey = func(space []byte) []byte {
			tmp := iter.tmp[:keySize]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipMap", err.Error())
				return nil
			}
			space = append(space, tmp...)
			return space
		}
	} else {
		switch keyType {
		case protocol.STRING:
			skipKey = iter.SkipBinary
		default:
			panic("unsupported type")
		}
	}
	if elemSize != 0 {
		skipElem = func(space []byte) []byte {
			tmp := iter.tmp[:elemSize]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipMap", err.Error())
				return nil
			}
			space = append(space, tmp...)
			return space
		}
	} else {
		switch elemType {
		case protocol.STRING:
			skipElem = iter.SkipBinary
		case protocol.LIST:
			skipElem = iter.SkipList
		case protocol.STRUCT:
			skipElem = iter.SkipStruct
		case protocol.MAP:
			skipElem = iter.SkipMap
		default:
			panic("unsupported type")
		}
	}
	for i := 0; i < length; i++ {
		space = skipKey(space)
		space = skipElem(space)
	}
	return space
}

func (iter *Iterator) ReadBool() bool {
	tmp := iter.tmp[:1]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBool", err.Error())
		return false
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadBool()
}

func (iter *Iterator) ReadInt8() int8 {
	return int8(iter.ReadUInt8())
}

func (iter *Iterator) ReadUInt8() uint8 {
	tmp := iter.tmp[:1]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUInt8", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUInt8()
}

func (iter *Iterator) ReadInt16() int16 {
	return int16(iter.ReadUInt16())
}

func (iter *Iterator) ReadUInt16() uint16 {
	tmp := iter.tmp[:2]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUInt16", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUInt16()
}

func (iter *Iterator) ReadInt32() int32 {
	return int32(iter.ReadUInt32())
}

func (iter *Iterator) ReadUInt32() uint32 {
	tmp := iter.tmp[:4]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUInt32", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUInt32()
}

func (iter *Iterator) ReadInt64() int64 {
	tmp := iter.allocate(8)
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadInt64", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadInt64()
}

func (iter *Iterator) ReadUInt64() uint64 {
	tmp := iter.tmp[:8]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUInt64", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadUInt64()
}

func (iter *Iterator) ReadFloat64() float64 {
	tmp := iter.tmp[:8]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadFloat64", err.Error())
		return 0
	}
	iter.real.Reset(nil, tmp)
	return iter.real.ReadFloat64()
}

func (iter *Iterator) ReadString() string {
	size := iter.ReadUInt32()
	tmp := iter.allocate(int(size))
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBinary", err.Error())
		return ""
	}
	return string(tmp)
}

func (iter *Iterator) ReadBinary() []byte {
	size := iter.ReadUInt32()
	tmp := iter.allocate(int(size))
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBinary", err.Error())
		return nil
	}
	return tmp
}

func (iter *Iterator) SkipBinary(space []byte) []byte {
	tmp := iter.tmp[:4]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("SkipBinary", err.Error())
		return nil
	}
	space = append(space, tmp...)
	size := uint32(tmp[3]) | uint32(tmp[2])<<8 | uint32(tmp[1])<<16 | uint32(tmp[0])<<24
	tmp = iter.allocate(int(size))
	_, err = io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("SkipBinary", err.Error())
		return nil
	}
	space = append(space, tmp...)
	return space
}

func (iter *Iterator) Read(ttype protocol.TType) interface{} {
	switch ttype {
	case protocol.BOOL:
		return iter.ReadBool()
	case protocol.I08:
		return iter.ReadInt8()
	case protocol.I16:
		return iter.ReadInt16()
	case protocol.I32:
		return iter.ReadInt32()
	case protocol.I64:
		return iter.ReadInt64()
	case protocol.DOUBLE:
		return iter.ReadFloat64()
	case protocol.STRING:
		return iter.ReadString()
	case protocol.LIST:
		return iter.ReadList()
	case protocol.MAP:
		return iter.ReadMap()
	case protocol.STRUCT:
		return iter.ReadStruct()
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) ReaderOf(ttype protocol.TType) func() interface{} {
	switch ttype {
	case protocol.BOOL:
		return func() interface{} {
			return iter.ReadBool()
		}
	case protocol.I08:
		return func() interface{} {
			return iter.ReadInt8()
		}
	case protocol.I16:
		return func() interface{} {
			return iter.ReadInt16()
		}
	case protocol.I32:
		return func() interface{} {
			return iter.ReadInt32()
		}
	case protocol.I64:
		return func() interface{} {
			return iter.ReadInt64()
		}
	case protocol.DOUBLE:
		return func() interface{} {
			return iter.ReadFloat64()
		}
	case protocol.STRING:
		return func() interface{} {
			return iter.ReadString()
		}
	case protocol.LIST:
		return func() interface{} {
			return iter.ReadList()
		}
	case protocol.MAP:
		return func() interface{} {
			return iter.ReadMap()
		}
	case protocol.STRUCT:
		return func() interface{} {
			return iter.ReadStruct()
		}
	default:
		panic("unsupported type")
	}
}

func getTypeSize(elemType protocol.TType) int {
	switch elemType {
	case protocol.BOOL, protocol.I08:
		return 1
	case protocol.I16:
		return 2
	case protocol.I32:
		return 4
	case protocol.I64, protocol.DOUBLE:
		return 8
	}
	return 0
}