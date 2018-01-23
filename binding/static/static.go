package static

import (
	"github.com/thrift-iterator/go/spi"
	"reflect"
)

type CodegenExtension struct {
	spi.Extension
	ExtTypes []reflect.Type
}

func (ext *CodegenExtension) MangledName() string {
	// TODO: hash extension to represent different config
	return "default"
}