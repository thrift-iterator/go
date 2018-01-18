package binary

import "github.com/thrift-iterator/go/protocol"

func (iter *Iterator) Discard(ttype protocol.TType) {
	switch ttype {
	case protocol.BOOL, protocol.I08:
		iter.buf = iter.buf[1:]
	case protocol.I16:
		iter.buf = iter.buf[2:]
	case protocol.I32:
		iter.buf = iter.buf[4:]
	case protocol.I64, protocol.DOUBLE:
		iter.buf = iter.buf[8:]
	case protocol.STRING:
		b := iter.buf
		size := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
		iter.buf = iter.buf[size:]
	case protocol.LIST:
		iter.SkipList(nil)
	case protocol.MAP:
		iter.SkipMap(nil)
	case protocol.STRUCT:
		iter.SkipStruct(nil)
	default:
		panic("unsupported type")
	}
}