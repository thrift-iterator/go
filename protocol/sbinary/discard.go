package sbinary

import (
	"github.com/thrift-iterator/go/protocol"
)

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
	case protocol.TypeBool, protocol.TypeI08:
		iter.readSmall(1)
	case protocol.TypeI16:
		iter.readSmall(2)
	case protocol.TypeI32:
		iter.readSmall(4)
	case protocol.TypeI64, protocol.TypeDouble:
		iter.readSmall(8)
	case protocol.TypeString:
		iter.discardBinary()
	case protocol.TypeList:
		iter.discardList()
	case protocol.TypeMap:
		iter.discardMap()
	case protocol.TypeStruct:
		iter.discardStruct()
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) discardMap() {
	tmp := iter.readSmall(6)
	iter.real.Reset(nil, tmp)
	keyType, elemType, length := iter.real.ReadMapHeader()
	keySize := getTypeSize(keyType)
	elemSize := getTypeSize(elemType)
	if keySize != 0 && elemSize != 0 {
		iter.readLarge(length * (keySize + elemSize))
		return
	}
	var discardKey func()
	var discardElem func()
	if keySize != 0 {
		discardKey = func() {
			iter.readSmall(keySize)
			return
		}
	} else {
		switch keyType {
		case protocol.TypeString:
			discardKey = iter.discardBinary
		default:
			panic("unsupported type")
		}
	}
	if elemSize != 0 {
		discardElem = func() {
			iter.readSmall(elemSize)
			return
		}
	} else {
		switch elemType {
		case protocol.TypeString:
			discardElem = iter.discardBinary
		case protocol.TypeList:
			discardElem = iter.discardList
		case protocol.TypeStruct:
			discardElem = iter.discardStruct
		case protocol.TypeMap:
			discardElem = iter.discardMap
		default:
			panic("unsupported type")
		}
	}
	for i := 0; i < length; i++ {
		discardKey()
		discardElem()
	}
	return
}

func (iter *Iterator) discardStruct() {
	for {
		tmp := iter.readSmall(1)
		fieldType := protocol.TType(tmp[0])
		switch fieldType {
		case protocol.TypeStop:
			return
		case protocol.TypeBool, protocol.TypeI08:
			iter.readSmall(3) // 1 + 2
		case protocol.TypeI16:
			iter.readSmall(4) // 2 + 2
		case protocol.TypeI32:
			iter.readSmall(6) // 4 + 2
		case protocol.TypeI64, protocol.TypeDouble:
			iter.readSmall(10) // 8 + 2
		case protocol.TypeList:
			iter.readSmall(2)
			iter.discardList()
		case protocol.TypeMap:
			iter.readSmall(2)
			iter.discardMap()
		case protocol.TypeString:
			iter.readSmall(2)
			iter.discardBinary()
		case protocol.TypeStruct:
			iter.readSmall(2)
			iter.discardStruct()
		default:
			panic("unsupported type")
		}
	}
}

func (iter *Iterator) discardList() {
	tmp := iter.readSmall(5)
	iter.real.Reset(nil, tmp)
	elemType, length := iter.real.ReadListHeader()
	switch elemType {
	case protocol.TypeStop:
	case protocol.TypeBool, protocol.TypeI08:
		iter.readLarge(length)
	case protocol.TypeI16:
		iter.readLarge(length * 2)
	case protocol.TypeI32:
		iter.readLarge(length * 4)
	case protocol.TypeI64, protocol.TypeDouble:
		iter.readLarge(length * 8)
	case protocol.TypeString:
		for i := 0; i < length; i++ {
			iter.discardBinary()
		}
	case protocol.TypeList:
		for i := 0; i < length; i++ {
			iter.discardList()
		}
	case protocol.TypeMap:
		for i := 0; i < length; i++ {
			iter.discardMap()
		}
	case protocol.TypeStruct:
		for i := 0; i < length; i++ {
			iter.discardStruct()
		}
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) discardBinary() {
	tmp := iter.readSmall(4)
	size := uint32(tmp[3]) | uint32(tmp[2])<<8 | uint32(tmp[1])<<16 | uint32(tmp[0])<<24
	iter.readLarge(int(size))
}
