package binary

import (
	"fmt"
	"github.com/thrift-iterator/go/protocol"
)

type Iterator struct {
	buf   []byte
	Error error
}

func NewIterator(buf []byte) *Iterator {
	return &Iterator{buf: buf}
}

func (iter *Iterator) ReadStructCB(cb func(fieldType protocol.TType, fieldId protocol.FieldId)) {
	for iter.buf[0] != 0 {
		fieldType := iter.buf[0]
		fieldId := uint16(iter.buf[2]) | uint16(iter.buf[1])<<8
		iter.buf = iter.buf[3:]
		cb(protocol.TType(fieldType), protocol.FieldId(fieldId))
	}
	iter.buf = iter.buf[1:]
}

func (iter *Iterator) ReadInt64() int64 {
	b := iter.buf
	value := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	iter.buf = iter.buf[8:]
	return int64(value)
}

func (iter *Iterator) ReportError(operation string, err string) {
	if iter.Error == nil {
		iter.Error = fmt.Errorf("%s: %s", operation, err)
	}
}
