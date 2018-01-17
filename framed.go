package thrifter

import (
	"io"
	"encoding/binary"
	"github.com/thrift-iterator/go/protocol"
	"errors"
)

type framedDecoder struct {
	reader io.Reader
	iter   BufferedIterator
	tmp    []byte
}

func (decoder *framedDecoder) Decode(obj interface{}) error {
	msg, _ := obj.(*protocol.Message)
	if msg == nil {
		return errors.New("can only unmarshal protocol.Message")
	}
	if len(decoder.tmp) < 4 {
		decoder.tmp = make([]byte, 64)
	}
	tmp := decoder.tmp[:4]
	_, err := io.ReadFull(decoder.reader, tmp)
	if err != nil {
		return err
	}
	size := int(binary.BigEndian.Uint32(tmp))
	if len(decoder.tmp) < size {
		decoder.tmp = make([]byte, size)
	}
	tmp = decoder.tmp[:size]
	_, err = io.ReadFull(decoder.reader, tmp)
	if err != nil {
		return err
	}
	decoder.iter.Reset(tmp)
	msgRead := decoder.iter.ReadMessage()
	msg.Set(&msgRead)
	return nil
}
