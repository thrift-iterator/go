package compact

import (
	"github.com/thrift-iterator/go/protocol"
)

func (iter *Iterator) SkipMessage(space []byte) []byte {
	panic("not implemented")
}

func (iter *Iterator) SkipStruct(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.skipStruct()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) skipStruct() {
	iter.ReadStructHeader()
	for {
		fieldType, _ := iter.ReadStructField()
		if fieldType == 0 {
			return
		}
		switch fieldType {
		case protocol.BOOL, protocol.I08:
			iter.ReadInt8()
		case protocol.I16:
			iter.ReadInt16()
		case protocol.I32:
			iter.ReadInt32()
		case protocol.I64:
			iter.ReadInt64()
		case protocol.DOUBLE:
			iter.ReadFloat64()
		case protocol.STRING:
			iter.skipBinary()
		case protocol.LIST:
			iter.skipList()
		case protocol.STRUCT:
			iter.skipStruct()
		case protocol.MAP:
			iter.skipMap()
		default:
			panic("unsupported type")
		}
	}
}

func (iter *Iterator) SkipList(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.skipList()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) skipList() {
	elemType, length := iter.ReadListHeader()
	skipElem := iter.howToSkip(elemType)
	for i := 0; i < length; i++ {
		skipElem()
	}
}

func (iter *Iterator) SkipMap(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.skipMap()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) skipMap() {
	keyType, elemType, length := iter.ReadMapHeader()
	skipKey := iter.howToSkip(keyType)
	skipElem := iter.howToSkip(elemType)
	for i := 0; i < length; i++ {
		skipKey()
		skipElem()
	}
}

func (iter *Iterator) SkipBinary(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.skipBinary()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) skipBinary() {
	iter.ReadBinary()
}

func (iter *Iterator) howToSkip(elemType protocol.TType) func() {
	switch elemType {
	case protocol.BOOL, protocol.I08:
		return func() {
			iter.ReadInt8()
		}
	case protocol.I16:
		return func() {
			iter.ReadInt16()
		}
	case protocol.I32:
		return func() {
			iter.ReadInt32()
		}
	case protocol.I64:
		return func() {
			iter.ReadInt64()
		}
	case protocol.DOUBLE:
		return func() {
			iter.ReadFloat64()
		}
	case protocol.STRING:
		return iter.skipBinary
	case protocol.LIST:
		return iter.skipList
	case protocol.STRUCT:
		return iter.skipStruct
	case protocol.MAP:
		return iter.skipMap
	}
	panic("unsupported type")
}
