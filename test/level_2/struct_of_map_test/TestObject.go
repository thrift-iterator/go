package struct_of_map_test

type TestObject struct {
	Field1 map[int32]int64 `thrift:",1"`
}