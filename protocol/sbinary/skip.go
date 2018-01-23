package sbinary

import (
	"io"
	"github.com/thrift-iterator/go/protocol"
)

func (iter *Iterator) Skip(ttype protocol.TType, space []byte) []byte {
	switch ttype {
	case protocol.TypeBool, protocol.TypeI08:
		return iter.skip(space, 1)
	case protocol.TypeI16:
		return iter.skip(space, 2)
	case protocol.TypeI32:
		return iter.skip(space, 4)
	case protocol.TypeI64, protocol.TypeDouble:
		return iter.skip(space, 8)
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
	space = iter.skip(space, 4)
	space = iter.SkipBinary(space)
	space = iter.skip(space, 4)
	return space
}

func (iter *Iterator) SkipMap(space []byte) []byte {
	if space == nil {
		iter.recorder = []byte{}
	} else {
		iter.recorder = space
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
		iter.recorder = space
	}
	iter.discardStruct()
	// reuse buffer next time, save it on iter.space
	iter.space = iter.recorder
	iter.recorder = nil
	return iter.space
}

func (iter *Iterator) SkipBinary(space []byte) []byte {
	if space == nil {
		iter.recorder = []byte{}
	} else {
		iter.recorder = space
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
		iter.recorder = space
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