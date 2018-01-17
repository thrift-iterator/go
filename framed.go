package thrifter

import (
	"io"
	"encoding/binary"
	"github.com/thrift-iterator/go/protocol"
	"errors"
)

type framedDecoder struct {
	reader io.Reader
	iter   Iterator
	tmp    []byte
}

type framedEncoder struct {
	writer io.Writer
	stream Stream
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
	decoder.iter.Reset(nil, tmp)
	msgRead := decoder.iter.ReadMessage()
	msg.Set(&msgRead)
	return nil
}

func (encoder *framedEncoder) Encode(obj interface{}) error {
	msg, isMsg := obj.(protocol.Message)
	if !isMsg {
		return errors.New("can only unmarshal protocol.Message")
	}
	encoder.stream.Reset(nil)
	encoder.stream.WriteMessage(msg)
	buf := encoder.stream.Buffer()
	size := len(buf)
	_, err := encoder.writer.Write([]byte{
		byte(size >> 24), byte(size >> 16), byte(size >> 8), byte(size),
	})
	if err != nil {
		return err
	}
	_, err = encoder.writer.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
