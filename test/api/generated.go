
package test
import "github.com/v2pro/wombat/generic"
import "github.com/thrift-iterator/go/test/api/binding_test"
import "github.com/thrift-iterator/go/protocol/binary"
func init() {
generic.RegisterExpandedFunc("Decode_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator",Decode_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator)}
func DecodeSimpleValue_DT_ptr_int64_ST_ptr_binary__Iterator(dst *int64,src *binary.Iterator){
*dst = int64(src.ReadInt64())
	
}
func DecodeAnything_DT_ptr_int64_ST_ptr_binary__Iterator(dst *int64,src *binary.Iterator){


DecodeSimpleValue_DT_ptr_int64_ST_ptr_binary__Iterator(dst, src)

}
func DecodeStruct_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator(dst *binding_test.TestObject,src *binary.Iterator){


	
	

src.ReadStructHeader()
for {
	fieldType, fieldId := src.ReadStructField()
	if fieldType == 0 {
		return
	}
	switch fieldId {
		
			case 1:
				DecodeAnything_DT_ptr_int64_ST_ptr_binary__Iterator(&dst.Field1, src)
		
		default:
			src.Discard(fieldType)
	}
}
}
func DecodeAnything_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator(dst *binding_test.TestObject,src *binary.Iterator){


DecodeStruct_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator(dst, src)

}
func Decode_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator(dst interface{},src interface{}){

DecodeAnything_DT_ptr_binding_test__TestObject_ST_ptr_binary__Iterator(dst.(*binding_test.TestObject), src.(*binary.Iterator))

}