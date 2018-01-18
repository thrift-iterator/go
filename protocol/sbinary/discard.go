package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
)

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
	case protocol.BOOL, protocol.I08:
		iter.discard(1)
	case protocol.I16:
		iter.discard(2)
	case protocol.I32:
		iter.discard(4)
	case protocol.I64, protocol.DOUBLE:
		iter.discard(8)
	case protocol.STRING:
		iter.discardBinary()
	case protocol.LIST:
		iter.discardList()
	case protocol.MAP:
		iter.discardMap()
	case protocol.STRUCT:
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
		case protocol.STRING:
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
		case protocol.STRING:
			discardElem = iter.discardBinary
		case protocol.LIST:
			discardElem = iter.discardList
		case protocol.STRUCT:
			discardElem = iter.discardStruct
		case protocol.MAP:
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
		case protocol.STOP:
			return
		case protocol.I64, protocol.DOUBLE:
			tmp := iter.tmp[:10]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
		case protocol.LIST:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardList()
		case protocol.MAP:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardMap()
		case protocol.STRING:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return
			}
			iter.discardBinary()
		case protocol.STRUCT:
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
	case protocol.STOP:
		return
	case protocol.I64, protocol.DOUBLE:
		tmp := iter.allocate(length * 8)
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("discardList", err.Error())
			return
		}
		return
	case protocol.STRING:
		for i := 0; i < length; i++ {
			iter.discardBinary()
		}
		return
	case protocol.LIST:
		for i := 0; i < length; i++ {
			iter.discardList()
		}
		return
	case protocol.MAP:
		for i := 0; i < length; i++ {
			iter.discardMap()
		}
		return
	case protocol.STRUCT:
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
