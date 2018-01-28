package thrifter

import (
	"github.com/thrift-iterator/go/general"
	"github.com/thrift-iterator/go/protocol"
	"io"
	"github.com/thrift-iterator/go/spi"
)

type Encoder struct {
	cfg    *frozenConfig
	stream spi.Stream
}

func (encoder *Encoder) EncodeMessage(msg general.Message) error {
	return encoder.Encode(msg)
}

func (encoder *Encoder) EncodeMessageHeader(msgHeader protocol.MessageHeader) error {
	return encoder.Encode(msgHeader)
}

func (encoder *Encoder) EncodeMessageArguments(msgArgs general.Struct) error {
	return encoder.Encode(msgArgs)
}

func (encoder *Encoder) Reset(writer io.Writer) {
	encoder.stream.Reset(writer)
}

func (encoder *Encoder) Buffer() []byte {
	return encoder.stream.Buffer()
}