package spi

import "github.com/thrift-iterator/go/protocol"

func DiscardList(iter Iterator) {
	elemType, size := iter.ReadListHeader()
	discardElem := howToDiscard(iter, elemType)
	for i := 0; i < size; i++ {
		discardElem()
	}
}

func DiscardStruct(iter Iterator) {
	iter.ReadStructHeader()
	for {
		fieldType, _ := iter.ReadStructField()
		if fieldType == 0 {
			return
		}
		iter.Discard(fieldType)
	}
}

func DiscardMap(iter Iterator) {
	keyType, elemType, size := iter.ReadMapHeader()
	discardKey := howToDiscard(iter, keyType)
	discardElem := howToDiscard(iter, elemType)
	for i := 0; i < size; i++ {
		discardKey()
		discardElem()
	}
}

func howToDiscard(iter Iterator, elemType protocol.TType) func() {
	switch elemType {
	case protocol.TypeBool, protocol.TypeI08,
		protocol.TypeI16, protocol.TypeI32, protocol.TypeI64,
		protocol.TypeDouble, protocol.TypeString:
		return func() { iter.Discard(elemType) }
	case protocol.TypeList:
		return func() { DiscardList(iter) }
	case protocol.TypeStruct:
		return func() { DiscardStruct(iter) }
	case protocol.TypeMap:
		return func() { DiscardMap(iter) }
	}
	panic("unsupported type")
}
