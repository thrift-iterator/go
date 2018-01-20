package struct_of_struct_test

type TestObject struct {
	Field1 EmbeddedObject `thrift:",1"`
}

type EmbeddedObject struct {
	Field1 string `thrift:",1"`
}