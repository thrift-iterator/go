package compact

import "github.com/thrift-iterator/go/protocol"

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
	case protocol.TypeBool, protocol.TypeI08:
		iter.ReadInt8()
	case protocol.TypeI16:
		iter.ReadInt16()
	case protocol.TypeI32:
		iter.ReadInt32()
	case protocol.TypeI64:
		iter.ReadInt64()
	case protocol.TypeDouble:
		iter.ReadFloat64()
	case protocol.TypeString:
		iter.discardBinary()
	case protocol.TypeList:
		iter.discardList()
	case protocol.TypeStruct:
		iter.discardStruct()
	case protocol.TypeMap:
		iter.discardMap()
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) discardMessageHeader() {
	iter.consume(2)
	iter.readVarInt32()
	iter.discardBinary()
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
	case protocol.TypeBool, protocol.TypeI08:
		return func() {
			iter.ReadInt8()
		}
	case protocol.TypeI16:
		return func() {
			iter.ReadInt16()
		}
	case protocol.TypeI32:
		return func() {
			iter.ReadInt32()
		}
	case protocol.TypeI64:
		return func() {
			iter.ReadInt64()
		}
	case protocol.TypeDouble:
		return func() {
			iter.ReadFloat64()
		}
	case protocol.TypeString:
		return iter.discardBinary
	case protocol.TypeList:
		return iter.discardList
	case protocol.TypeStruct:
		return iter.discardStruct
	case protocol.TypeMap:
		return iter.discardMap
	}
	panic("unsupported type")
}