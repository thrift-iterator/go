package binary

import "github.com/thrift-iterator/go/protocol"

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
	case protocol.TypeBool, protocol.TypeI08:
		iter.buf = iter.buf[1:]
	case protocol.TypeI16:
		iter.buf = iter.buf[2:]
	case protocol.TypeI32:
		iter.buf = iter.buf[4:]
	case protocol.TypeI64, protocol.TypeDouble:
		iter.buf = iter.buf[8:]
	case protocol.TypeString:
		b := iter.buf
		size := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
		iter.buf = iter.buf[size:]
	case protocol.TypeList:
		iter.SkipList(nil)
	case protocol.TypeMap:
		iter.SkipMap(nil)
	case protocol.TypeStruct:
		iter.SkipStruct(nil)
	default:
		panic("unsupported type")
	}
}