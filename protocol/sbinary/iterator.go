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

func NewIterator(reader io.Reader) *Iterator {
	return &Iterator{
		reader: reader, real: binary.NewIterator(nil),
		tmp:    make([]byte, 256),
		space:  make([]byte, 256),
	}
}

func (iter *Iterator) allocate(nBytes int) []byte {
	if len(iter.tmp) < nBytes {
		iter.tmp = make([]byte, nBytes)
	}
	return iter.tmp[:nBytes]
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
	panic("not implemented")
}
func (iter *Iterator) ReadMessage() protocol.Message {
	panic("not implemented")
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

func (iter *Iterator) ReadStructField() (fieldType protocol.TType, fieldId protocol.FieldId) {
	tmp := iter.tmp[:3]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadStructField", err.Error())
		return protocol.STOP, 0
	}
	iter.real.Reset(tmp)
	return iter.real.ReadStructField()
}

func (iter *Iterator) ReadStruct() map[protocol.FieldId]interface{} {
	buf := iter.SkipStruct(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(buf)
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
	iter.real.Reset(tmp)
	return iter.real.ReadListHeader()
}

func (iter *Iterator) ReadList() []interface{} {
	buf := iter.SkipList(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(buf)
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
	iter.real.Reset(tmp)
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
	iter.real.Reset(tmp)
	return iter.real.ReadMapHeader()
}

func (iter *Iterator) ReadMap() map[interface{}]interface{} {
	buf := iter.SkipMap(iter.space[:0])
	if iter.err != nil {
		return nil
	}
	iter.space = buf
	iter.real.Reset(buf)
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
	iter.real.Reset(tmp)
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
	panic("unsupported type")
}

func (iter *Iterator) ReadBool() bool {
	tmp := iter.tmp[:1]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadBool", err.Error())
		return false
	}
	iter.real.Reset(tmp)
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
	iter.real.Reset(tmp)
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
	iter.real.Reset(tmp)
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
	iter.real.Reset(tmp)
	return iter.real.ReadUInt32()
}

func (iter *Iterator) ReadInt64() int64 {
	tmp := iter.allocate(8)
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadInt64", err.Error())
		return 0
	}
	iter.real.Reset(tmp)
	return iter.real.ReadInt64()
}

func (iter *Iterator) ReadUInt64() uint64 {
	tmp := iter.tmp[:8]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadUInt64", err.Error())
		return 0
	}
	iter.real.Reset(tmp)
	return iter.real.ReadUInt64()
}

func (iter *Iterator) ReadFloat64() float64 {
	tmp := iter.tmp[:8]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("ReadFloat64", err.Error())
		return 0
	}
	iter.real.Reset(tmp)
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
	panic("not implemented")
}
func (iter *Iterator) Read(ttype protocol.TType) interface{} {
	panic("not implemented")
}
func (iter *Iterator) ReaderOf(ttype protocol.TType) func() interface{} {
	panic("not implemented")
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