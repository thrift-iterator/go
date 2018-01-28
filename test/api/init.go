package test

import (
	"github.com/v2pro/wombat/generic"
	"github.com/thrift-iterator/go"
	"github.com/thrift-iterator/go/test/api/binding_test"
)

var api = thrifter.Config{
	Protocol: thrifter.ProtocolBinary,
}.Froze()

//go:generate go install github.com/thrift-iterator/go/cmd/thrifter
//go:generate $GOPATH/bin/thrifter -pkg github.com/thrift-iterator/go/test/api
func init() {
	generic.Declare(func() {
		api.WillDecodeFromBuffer(
			(*binding_test.TestObject)(nil),
		)
	})
}
