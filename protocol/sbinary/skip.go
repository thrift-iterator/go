package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
)


func (iter *Iterator) SkipMessage(space []byte) []byte {
	space = iter.skip(space, 4)
	space = iter.SkipBinary(space)
	space = iter.skip(space, 4)
	space = iter.SkipStruct(space)
	return space
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
		case protocol.TypeString:
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
		case protocol.TypeString:
			skipElem = iter.SkipBinary
		case protocol.TypeList:
			skipElem = iter.SkipList
		case protocol.TypeStruct:
			skipElem = iter.SkipStruct
		case protocol.TypeMap:
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
		case protocol.TypeStop:
			return space
		case protocol.TypeI64, protocol.TypeDouble:
			tmp := iter.tmp[:10]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
		case protocol.TypeList:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipList(space)
		case protocol.TypeMap:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipMap(space)
		case protocol.TypeString:
			tmp := iter.tmp[:2]
			_, err := io.ReadFull(iter.reader, tmp)
			if err != nil {
				iter.ReportError("SkipStruct", err.Error())
				return nil
			}
			space = append(space, tmp...)
			space = iter.SkipBinary(space)
		case protocol.TypeStruct:
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
	case protocol.TypeStop:
		return nil
	case protocol.TypeI64, protocol.TypeDouble:
		tmp := iter.allocate(length * 8)
		_, err := io.ReadFull(iter.reader, tmp)
		if err != nil {
			iter.ReportError("SkipList", err.Error())
			return nil
		}
		space = append(space, tmp...)
		return space
	case protocol.TypeString:
		for i := 0; i < length; i++ {
			space = iter.SkipBinary(space)
		}
		return space
	case protocol.TypeList:
		for i := 0; i < length; i++ {
			space = iter.SkipList(space)
		}
		return space
	case protocol.TypeMap:
		for i := 0; i < length; i++ {
			space = iter.SkipMap(space)
		}
		return space
	case protocol.TypeStruct:
		for i := 0; i < length; i++ {
			space = iter.SkipStruct(space)
		}
		return space
	default:
		panic("unsupported type")
	}
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