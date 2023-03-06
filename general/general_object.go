package general

import (
	"bytes"
	"encoding/json"
	"github.com/thrift-iterator/go/protocol"
)

var (
	_ json.Marshaler = (*List)(nil)
	_ json.Marshaler = (*Map)(nil)
	_ json.Marshaler = (*Struct)(nil)
)

type Object interface {
	Get(path ...interface{}) interface{}
}

type List []interface{}

func (obj List) Get(path ...interface{}) interface{} {
	if len(path) == 0 {
		return obj
	}
	elem := obj[path[0].(int)]
	if len(path) == 1 {
		return elem
	}
	return elem.(Object).Get(path[1:]...)
}

func (l *List) MarshalJSON() ([]byte, error) {
	return json.Marshal((*[]interface{})(l))
}

type Map map[interface{}]interface{}

func (obj Map) Get(path ...interface{}) interface{} {
	if len(path) == 0 {
		return obj
	}
	elem := obj[path[0]]
	if len(path) == 1 {
		return elem
	}
	return elem.(Object).Get(path[1:]...)
}

func (m Map) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
		return []byte("{}"), nil
	}
	buf := bytes.NewBuffer([]byte("{"))
	for k, v := range m {
		buf.WriteString(`"`)
		buf.WriteString(k.(string))
		buf.WriteString(`":`)
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
		buf.WriteString(",")
	}
	buf.Truncate(buf.Len() - 1)
	buf.WriteString("}")
	return buf.Bytes(), nil
}

type Struct map[protocol.FieldId]interface{}

func (obj Struct) Get(path ...interface{}) interface{} {
	if len(path) == 0 {
		return obj
	}
	elem := obj[path[0].(protocol.FieldId)]
	if len(path) == 1 {
		return elem
	}
	return elem.(Object).Get(path[1:]...)
}

func (s *Struct) MarshalJSON() ([]byte, error) {
	return json.Marshal((*map[protocol.FieldId]interface{})(s))
}

type Message struct {
	protocol.MessageHeader
	Arguments Struct
}
