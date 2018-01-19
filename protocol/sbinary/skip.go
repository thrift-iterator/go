package sbinary

import (
	"io"
)


func (iter *Iterator) SkipMessage(space []byte) []byte {
	space = iter.skip(space, 4)
	space = iter.SkipBinary(space)
	space = iter.skip(space, 4)
	space = iter.SkipStruct(space)
	return space
}

func (iter *Iterator) SkipMap(space []byte) []byte {
	if space == nil {
		iter.recorder = []byte{}
	} else {
		iter.recorder = space[:0]
	}
	iter.discardMap()
	iter.space = iter.recorder
	iter.recorder = nil
	return iter.space
}

func (iter *Iterator) SkipStruct(space []byte) []byte {
	if space == nil {
		iter.recorder = []byte{}
	} else {
		iter.recorder = space[:0]
	}
	iter.discardStruct()
	iter.space = iter.recorder
	iter.recorder = nil
	return iter.space
}

func (iter *Iterator) SkipBinary(space []byte) []byte {
	if space == nil {
		iter.recorder = []byte{}
	} else {
		iter.recorder = space[:0]
	}
	iter.discardBinary()
	iter.space = iter.recorder
	iter.recorder = nil
	return iter.space
}

func (iter *Iterator) SkipList(space []byte) []byte {
	if space == nil {
		iter.recorder = []byte{}
	} else {
		iter.recorder = space[:0]
	}
	iter.discardList()
	iter.space = iter.recorder
	iter.recorder = nil
	return iter.space
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