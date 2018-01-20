package struct_of_pointer_test

type StructOf2Ptr struct {
	Field1 *int `thrift:",1"`
	Field2 *int `thrift:",2"`
}