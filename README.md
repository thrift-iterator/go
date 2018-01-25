# thrifter

decode/encode thrift message without IDL

Why?

* because IDL generated model is ugly and inflexible, it is seldom used in application directly. instead we define another model, which leads to bad performance.
  * bytes need to be copied twice 
  * more objects to gc
* thrift proxy can not know all possible IDL in advance, we need to decode/encode in a generic way to modify embedded header.
* official thrift library for go is slow, verified in several benchmarks. It is even slower than [json-iterator](https://github.com/json-iterator/go)

# works like encoding/json

`encoding/json` has a super simple api to encode/decode json.
thrifter mimic the same api.

```go
import "github.com/thrift-iterator/go"
// marshal to thrift
thriftEncodedBytes, err := thrifter.Marshal([]int{1, 2, 3})
// unmarshal back
var val []int
err = thrifter.Unmarshal(thriftEncodedBytes, &val)
```

event struct data binding is supported

```go
import "github.com/thrift-iterator/go"

type NewOrderRequest struct {
    Lines []NewOrderLine `thrift:",1"`
}

type NewOrderLine struct {
    ProductId string `thrift:",1"`
    Quantity int `thrift:",2"`
}

// marshal to thrift
thriftEncodedBytes, err := thrifter.Marshal(NewOrderRequest{
	Lines: []NewOrderLine{
		{"apple", 1},
		{"orange", 2},
	}
})
// unmarshal back
var val NewOrderRequest
err = thrifter.Unmarshal(thriftEncodedBytes, &val)
```

# without IDL

you do not need to define IDL. you do not need to use static code generation.
you do not event need to define struct.

```go
import "github.com/thrift-iterator/go"
import "github.com/thrift-iterator/go/general"

// msg is of type general.Message
msg, err := thrifter.UnmarshalMessage(thriftEncodedBytes)
// the RPC call method name, type is string
fmt.Println(msg.MessageName)
// the RPC call arguments, type is general.Struct
fmt.Println(msg.MessageArgs)
```

what is `general.Struct`, it is defined as a map

```go
type FieldId int16
type Struct map[FieldId]interface{}
```

we can extract out specific argument using one line

```go
productId := msg.MessageArgs.Get(
	protocol.FieldId(1), // lines of request
	0, // the first line
	protocol.FieldId(1), // product id
).(string)
```

You can unmarshal any thrift bytes into general objects. And you can marshal them back.

# Partial decoding

TODO