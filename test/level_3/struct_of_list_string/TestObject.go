package struct_of_list_string

type TestObject struct {
	Field1 []string `thrift:",1"`
	Field2 int64    `thrift:",2"`
}
