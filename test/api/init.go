package test

import (
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test/api/binding_test"
)

//go:generate go install github.com/v2pro/wombat/cmd/wombat-codegen
//go:generate $GOPATH/bin/wombat-codegen -pkg github.com/thrift-iterator/go/test/api
func init() {
	generic.Declare(func() {
		thrifter.DefaultConfig.WillDecodeFromBuffer(
			(*binding_test.TestObject)(nil),
		)
	})
}
