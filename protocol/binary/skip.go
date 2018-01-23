package binary

import "github.com/thrift-iterator/go/protocol"


func (iter *Iterator) Skip(ttype protocol.TType, space []byte) []byte {
	switch ttype {
	case protocol.TypeBool, protocol.TypeI08:
		skipped := iter.buf[:1]
		iter.buf = iter.buf[1:]
		if len(space) > 0 {
			space = append(space, skipped...)
			return space
		}
		return skipped
	case protocol.TypeI16:
		skipped := iter.buf[:2]
		iter.buf = iter.buf[2:]
		if len(space) > 0 {
			space = append(space, skipped...)
			return space
		}
		return skipped
	case protocol.TypeI32:
		skipped := iter.buf[:4]
		iter.buf = iter.buf[4:]
		if len(space) > 0 {
			space = append(space, skipped...)
			return space
		}
		return skipped
	case protocol.TypeI64, protocol.TypeDouble:
		skipped := iter.buf[:8]
		iter.buf = iter.buf[8:]
		if len(space) > 0 {
			space = append(space, skipped...)
			return space
		}
		return skipped
	case protocol.TypeString:
		return iter.SkipBinary(space)
	case protocol.TypeList:
		return iter.SkipList(space)
	case protocol.TypeMap:
		return iter.SkipMap(space)
	case protocol.TypeStruct:
		return iter.SkipStruct(space)
	default:
		panic("unsupported type")
	}
}

func (iter *Iterator) SkipMessageHeader(space []byte) []byte {
	bufBeforeSkip := iter.buf
	iter.buf = iter.buf[4:]
	skippedBytes := 4 + len(iter.skipBinary()) + 4
	skipped := bufBeforeSkip[:skippedBytes]
	iter.buf = bufBeforeSkip[skippedBytes:]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) SkipStruct(space []byte) []byte {
	bufBeforeSkip := iter.buf
	skippedBytes := 0
	for {
		fieldType := protocol.TType(iter.buf[0])
		if fieldType == 0 {
			iter.buf = iter.buf[1:]
			skippedBytes += 1
			if len(space) > 0 {
				return append(space, bufBeforeSkip[:skippedBytes]...)
			}
			return bufBeforeSkip[:skippedBytes]
		}
		switch fieldType {
		case protocol.TypeBool, protocol.TypeI08:
			iter.buf = iter.buf[4:]
			skippedBytes += 4
		case protocol.TypeI16:
			iter.buf = iter.buf[5:]
			skippedBytes += 5
		case protocol.TypeI32:
			iter.buf = iter.buf[7:]
			skippedBytes += 7
		case protocol.TypeI64, protocol.TypeDouble:
			iter.buf = iter.buf[11:]
			skippedBytes += 11
		case protocol.TypeString:
			b := iter.buf
			size := uint32(b[6]) | uint32(b[5])<<8 | uint32(b[4])<<16 | uint32(b[3])<<24
			skippedBytes += int(size)
			skippedBytes += 7
			iter.buf = bufBeforeSkip[skippedBytes:]
		case protocol.TypeList:
			iter.buf = iter.buf[3:]
			skippedBytes += len(iter.SkipList(nil))
			skippedBytes += 3
		case protocol.TypeMap:
			iter.buf = iter.buf[3:]
			skippedBytes += len(iter.SkipMap(nil))
			skippedBytes += 3
		case protocol.TypeStruct:
			iter.buf = iter.buf[3:]
			skippedBytes += len(iter.SkipStruct(nil))
			skippedBytes += 3
		default:
			panic("unsupported type")
		}
	}
}

func (iter *Iterator) SkipList(space []byte) []byte {
	if len(space) > 0 {
		return append(space, iter.skipList()...)
	}
	return iter.skipList()
}

func (iter *Iterator) skipList() []byte {
	bufBeforeSkip := iter.buf
	elemType := protocol.TType(bufBeforeSkip[0])
	length := uint32(bufBeforeSkip[4]) | uint32(bufBeforeSkip[3])<<8 | uint32(bufBeforeSkip[2])<<16 | uint32(bufBeforeSkip[1])<<24
	switch elemType {
	case protocol.TypeBool, protocol.TypeI08:
		size := 5 + length
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.TypeI16:
		size := 5 + length*2
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.TypeI32:
		size := 5 + length*4
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.TypeI64, protocol.TypeDouble:
		size := 5 + length*8
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.TypeString:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.ReadBinary())
			skippedBytes += 4
		}
		iter.buf = bufBeforeSkip[skippedBytes:]
		return bufBeforeSkip[:skippedBytes]
	case protocol.TypeList:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipList(nil))
		}
		return bufBeforeSkip[:skippedBytes]
	case protocol.TypeMap:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipMap(nil))
		}
		return bufBeforeSkip[:skippedBytes]
	case protocol.TypeStruct:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipStruct(nil))
		}
		return bufBeforeSkip[:skippedBytes]
	}
	panic("unsupported type")
}

func (iter *Iterator) SkipMap(space []byte) []byte {
	if len(space) > 0 {
		return append(space, iter.skipMap()...)
	}
	return iter.skipMap()
}

func (iter *Iterator) skipMap() []byte {
	bufBeforeSkip := iter.buf
	keyType := protocol.TType(bufBeforeSkip[0])
	elemType := protocol.TType(bufBeforeSkip[1])
	length := uint32(bufBeforeSkip[5]) | uint32(bufBeforeSkip[4])<<8 | uint32(bufBeforeSkip[3])<<16 | uint32(bufBeforeSkip[2])<<24
	keySize := getTypeSize(keyType)
	elemSize := getTypeSize(elemType)
	if keySize != 0 && elemSize != 0 {
		size := 6 + int(length)*(elemSize+keySize)
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	}
	var skipKey func(space []byte) []byte
	var skipElem func(space []byte) []byte
	if keySize != 0 {
		skipKey = func(space []byte) []byte {
			skipped := iter.buf[:keySize]
			iter.buf = iter.buf[keySize:]
			return skipped
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
			skipped := iter.buf[:elemSize]
			iter.buf = iter.buf[elemSize:]
			return skipped
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
	skippedBytes := 6
	iter.buf = iter.buf[6:]
	for i := uint32(0); i < length; i++ {
		skippedBytes += len(skipKey(nil))
		skippedBytes += len(skipElem(nil))
	}
	return bufBeforeSkip[:skippedBytes]
}

func (iter *Iterator) SkipBinary(space []byte) []byte {
	if len(space) > 0 {
		return append(space, iter.skipBinary()...)
	}
	return iter.skipBinary()
}

func (iter *Iterator) skipBinary() []byte {
	b := iter.buf
	size := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	skipped := iter.buf[:4+size]
	iter.buf = iter.buf[4+size:]
	return skipped
}

func getTypeSize(elemType protocol.TType) int {
	switch elemType {
	case protocol.TypeBool, protocol.TypeI08:
		return 1
	case protocol.TypeI16:
		return 2
	case protocol.TypeI32:
		return 4
	case protocol.TypeI64, protocol.TypeDouble:
		return 8
	}
	return 0
}
