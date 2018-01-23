package general

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
	"github.com/thrift-iterator/go/protocol"
)

type Extension struct {
}

func (ext *Extension) EncoderOf(valType reflect.Type) spi.ValEncoder {
	switch valType {
	case reflect.TypeOf(([]interface{})(nil)):
		return &generalListEncoder{}
	case reflect.TypeOf((map[interface{}]interface{})(nil)):
		return &generalMapEncoder{}
	case reflect.TypeOf((map[protocol.FieldId]interface{})(nil)):
		return &generalStructEncoder{}
	case reflect.TypeOf((*protocol.Message)(nil)).Elem():
		return &messageEncoder{}
	case reflect.TypeOf((*protocol.MessageHeader)(nil)).Elem():
		return &messageHeaderEncoder{}
	}
	return nil
}

func (ext *Extension) DecoderOf(valType reflect.Type) spi.ValDecoder {
	switch valType {
	case reflect.TypeOf((*[]interface{})(nil)):
		return &generalListDecoder{}
	case reflect.TypeOf((*map[interface{}]interface{})(nil)):
		return &generalMapDecoder{}
	case reflect.TypeOf((*map[protocol.FieldId]interface{})(nil)):
		return &generalStructDecoder{}
	case reflect.TypeOf((*protocol.Message)(nil)):
		return &messageDecoder{}
	case reflect.TypeOf((*protocol.MessageHeader)(nil)):
		return &messageHeaderDecoder{}
	}
	return nil
}
