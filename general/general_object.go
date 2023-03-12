package general

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/thrift-iterator/go/protocol"
	"reflect"
	"strconv"
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
		i, errKey := writeKey(buf, k)
		if errKey != nil {
			return i, errKey
		}
		b, errValue := json.Marshal(v)
		if errValue != nil {
			return nil, errValue
		}
		buf.Write(b)
		buf.WriteString(",")
	}
	buf.Truncate(buf.Len() - 1)
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func writeKey(buf *bytes.Buffer, k interface{}) ([]byte, error) {
	buf.WriteString(`"`)
	switch k.(type) {
	case string:
		buf.WriteString(k.(string))
	case byte:
		buf.WriteString(strconv.FormatInt(int64(k.(byte)), 10))
	case int64:
		buf.WriteString(strconv.FormatInt(k.(int64), 10))
	case int32:
		buf.WriteString(strconv.FormatInt(int64(k.(int32)), 10))
	case int:
		buf.WriteString(strconv.FormatInt(int64(k.(int)), 10))
	case int16:
		buf.WriteString(strconv.FormatInt(int64(k.(int16)), 10))
	case float64:
		buf.WriteString(strconv.FormatFloat(k.(float64), 'f', -1, 64))
	case bool:
		buf.WriteString(strconv.FormatBool(k.(bool)))
	default:
		return nil, errors.New("unsupported map key type " + reflect.TypeOf(k).String())
	}
	buf.WriteString(`":`)
	return nil, nil
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
