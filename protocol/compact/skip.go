package compact

import "github.com/thrift-iterator/go/protocol"

func (iter *Iterator) SkipMessage(space []byte) []byte {
	panic("not implemented")
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
		case protocol.BOOL, protocol.I08:
			iter.buf = iter.buf[4:]
			skippedBytes += 4
		case protocol.I16:
			iter.buf = iter.buf[5:]
			skippedBytes += 5
		case protocol.I32:
			iter.buf = iter.buf[7:]
			skippedBytes += 7
		case protocol.I64, protocol.DOUBLE:
			iter.buf = iter.buf[11:]
			skippedBytes += 11
		case protocol.STRING:
			b := iter.buf
			size := uint32(b[6]) | uint32(b[5])<<8 | uint32(b[4])<<16 | uint32(b[3])<<24
			skippedBytes += int(size)
			skippedBytes += 7
			iter.buf = bufBeforeSkip[skippedBytes:]
		case protocol.LIST:
			iter.buf = iter.buf[3:]
			skippedBytes += len(iter.SkipList(nil))
			skippedBytes += 3
		case protocol.MAP:
			iter.buf = iter.buf[3:]
			skippedBytes += len(iter.SkipMap(nil))
			skippedBytes += 3
		case protocol.STRUCT:
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
	case protocol.BOOL, protocol.I08:
		size := 5 + length
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.I16:
		size := 5 + length*2
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.I32:
		size := 5 + length*4
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.I64, protocol.DOUBLE:
		size := 5 + length*8
		skipped := bufBeforeSkip[:size]
		iter.buf = bufBeforeSkip[size:]
		return skipped
	case protocol.STRING:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.ReadBinary())
			skippedBytes += 4
		}
		iter.buf = bufBeforeSkip[skippedBytes:]
		return bufBeforeSkip[:skippedBytes]
	case protocol.LIST:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipList(nil))
		}
		return bufBeforeSkip[:skippedBytes]
	case protocol.MAP:
		skippedBytes := 5
		iter.buf = iter.buf[5:]
		for i := uint32(0); i < length; i++ {
			skippedBytes += len(iter.SkipMap(nil))
		}
		return bufBeforeSkip[:skippedBytes]
	case protocol.STRUCT:
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
		case protocol.STRING:
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
		case protocol.STRING:
			skipElem = iter.SkipBinary
		case protocol.LIST:
			skipElem = iter.SkipList
		case protocol.STRUCT:
			skipElem = iter.SkipStruct
		case protocol.MAP:
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
	case protocol.BOOL, protocol.I08:
		return 1
	case protocol.I16:
		return 2
	case protocol.I32:
		return 4
	case protocol.I64, protocol.DOUBLE:
		return 8
	}
	return 0
}
