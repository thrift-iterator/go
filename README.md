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

