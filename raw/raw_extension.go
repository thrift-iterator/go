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
	case reflect.TypeOf((*Map)(nil)):
		return &rawMapDecoder{}
	}
	return nil
}

func (extension *Extension) EncoderOf(valType reflect.Type) spi.ValEncoder {
	switch valType {
	case reflect.TypeOf((*List)(nil)).Elem():
		return &rawListEncoder{}
	case reflect.TypeOf((*Map)(nil)).Elem():
		return &rawMapEncoder{}
	}
	return nil
}