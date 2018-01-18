package compact

import "github.com/thrift-iterator/go/protocol"

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
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
		iter.discardBinary()
	case protocol.LIST:
		iter.discardList()
	case protocol.STRUCT:
		iter.discardStruct()
	case protocol.MAP:
		iter.discardMap()
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) discardStruct() {
	iter.ReadStructHeader()
	for {
		fieldType, _ := iter.ReadStructField()
		if fieldType == 0 {
			return
		}
		iter.Discard(fieldType)
	}
}

func (iter *Iterator) discardList() {
	elemType, length := iter.ReadListHeader()
	discardElem := iter.howToDiscard(elemType)
	for i := 0; i < length; i++ {
		discardElem()
	}
}

func (iter *Iterator) discardMap() {
	keyType, elemType, length := iter.ReadMapHeader()
	discardKey := iter.howToDiscard(keyType)
	discardElem := iter.howToDiscard(elemType)
	for i := 0; i < length; i++ {
		discardKey()
		discardElem()
	}
}

func (iter *Iterator) discardBinary() {
	iter.ReadBinary()
}

func (iter *Iterator) howToDiscard(elemType protocol.TType) func() {
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
		return iter.discardBinary
	case protocol.LIST:
		return iter.discardList
	case protocol.STRUCT:
		return iter.discardStruct
	case protocol.MAP:
		return iter.discardMap
	}
	panic("unsupported type")
}