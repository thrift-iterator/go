package compact

import "github.com/thrift-iterator/go/protocol"

func (iter *Iterator) Skip(ttype protocol.TType, space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.Discard(ttype)
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) SkipMessageHeader(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.discardMessageHeader()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) SkipStruct(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.discardStruct()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) SkipList(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.discardList()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) SkipMap(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.discardMap()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}

func (iter *Iterator) SkipBinary(space []byte) []byte {
	bufBeforeSkip := iter.buf
	consumedBeforeSkip := iter.consumed
	iter.discardBinary()
	skipped := bufBeforeSkip[:iter.consumed-consumedBeforeSkip]
	if len(space) > 0 {
		return append(space, skipped...)
	}
	return skipped
}
