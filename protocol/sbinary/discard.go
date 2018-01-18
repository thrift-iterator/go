package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
)

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
	case protocol.TypeBool, protocol.TypeI08:
		iter.discard(1)
	case protocol.TypeI16:
		iter.discard(2)
	case protocol.TypeI32:
		iter.discard(4)
	case protocol.TypeI64, protocol.TypeDouble:
		iter.discard(8)
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
	tmp := iter.tmp[:6]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("discardMap", err.Error())
		return
	}
	iter.real.Reset(nil, tmp)
	keyType, elemType, length := iter.real.ReadMapHeader()
	keySize := getTypeSize(keyType)
	elemSize := getTypeSize(elemType)
	if keySize != 0 && elemSize != 0 {
		tmp := iter.allocate(length * (keySize + elemSize))
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("discardMap", err.Error())
			return
		}
		return
	}
	var discardKey func()
	var discardElem func()
	if keySize != 0 {
		discardKey = func() {
			tmp := iter.tmp[:keySize]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("discardMap", err.Error())
				return
			}
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
			tmp := iter.tmp[:elemSize]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("discardMap", err.Error())
				return
			}
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
		tmp := iter.tmp[:1]
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("SkipStruct", err.Error())
			return
		}
		fieldType := protocol.TType(tmp[0])
		switch fieldType {
		case protocol.TypeStop:
			return
		case protocol.TypeI64, protocol.TypeDouble:
			tmp := iter.tmp[:10]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
		case protocol.TypeList:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardList()
		case protocol.TypeMap:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardMap()
		case protocol.TypeString:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardBinary()
		case protocol.TypeStruct:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardStruct()
		default:
			panic("unsupported type")
		}
	}
}

func (iter *Iterator) discardList() {
	tmp := iter.tmp[:5]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("discardList", err.Error())
		return
	}
	iter.real.Reset(nil, tmp)
	elemType, length := iter.real.ReadListHeader()
	switch elemType {
	case protocol.TypeStop:
		return
	case protocol.TypeI64, protocol.TypeDouble:
		tmp := iter.allocate(length * 8)
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("discardList", err.Error())
			return
		}
		return
	case protocol.TypeString:
		for i := 0; i < length; i++ {
			iter.discardBinary()
		}
		return
	case protocol.TypeList:
		for i := 0; i < length; i++ {
			iter.discardList()
		}
		return
	case protocol.TypeMap:
		for i := 0; i < length; i++ {
			iter.discardMap()
		}
		return
	case protocol.TypeStruct:
		for i := 0; i < length; i++ {
			iter.discardStruct()
		}
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) discard(nBytes int) {
	tmp := iter.tmp[:nBytes]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("skip", err.Error())
	}
}

func (iter *Iterator) discardBinary() {
	tmp := iter.tmp[:4]
	_, err := io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("discardBinary", err.Error())
		return
	}
	size := uint32(tmp[3]) | uint32(tmp[2])<<8 | uint32(tmp[1])<<16 | uint32(tmp[0])<<24
	tmp = iter.allocate(int(size))
	_, err = io.ReadFull(iter.reader, tmp)
	if err != nil {
		iter.ReportError("discardBinary", err.Error())
		return
	}
}
