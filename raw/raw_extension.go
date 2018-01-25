package raw

import (
	"reflect"
	"github.com/thrift-iterator/go/spi"
)

type Extension struct {
}

func (extension *Extension) DecoderOf(valType reflect.Type) spi.ValDecoder {
	switch valType {
	case reflect.TypeOf((*List)(nil)):
		return &rawListDecoder{}
	}
	return nil
}

func (extension *Extension) EncoderOf(valType reflect.Type) spi.ValEncoder {
	switch valType {
	case reflect.TypeOf((*List)(nil)).Elem():
		return &rawListEncoder{}
	}
	return nil
}