package thrifter

import (
	"io"
)

type unframedDecoder struct {
	reader io.Reader
	iter   Iterator
	buf    []byte
	tmp    []byte
}

func (decoder *unframedDecoder) Decode(obj interface{}) error {
	decoder.readMessageHeader()
	return nil
}

func (decoder *unframedDecoder) readMessageHeader() {
}
