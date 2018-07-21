package struct_complex_test

type SubType struct {
	A int32 `thrift:"a,1"`
}

type Enum int32

const (
	Enum_A Enum = 1

	Enum_B Enum = 2
)

type Int int32

type TestObject struct {
	Av bool                             `thrift:"av,0"`
	Ap *bool                            `thrift:"ap,2,optional"`
	Bv int8                             `thrift:"bv,3"`
	Bp *int8                            `thrift:"bp,4,optional"`
	Cv int8                             `thrift:"cv,5"`
	Cp *int8                            `thrift:"cp,6,optional"`
	Dv int16                            `thrift:"dv,7"`
	Dp *int16                           `thrift:"dp,8,optional"`
	Ev int32                            `thrift:"ev,9"`
	Ep *int32                           `thrift:"ep,10,optional"`
	Fv int64                            `thrift:"fv,11"`
	Fp *int64                           `thrift:"fp,12,optional"`
	Gv float64                          `thrift:"gv,13"`
	Gp *float64                         `thrift:"gp,14,optional"`
	Hv string                           `thrift:"hv,15"`
	Hp *string                          `thrift:"hp,16,optional"`
	Iv []byte                           `thrift:"iv,17,optional"`
	Ip *[]byte                          `thrift:"ip,18,optional"`
	Jv []string                         `thrift:"jv,19,optional"`
	Jp *[]string                        `thrift:"jp,20,optional"`
	Kv map[string]bool                  `thrift:"kv,21,optional"`
	Kp *map[string]bool                 `thrift:"kp,22,optional"`
	Lv map[int32]SubType                `thrift:"lv,23,optional"`
	Lp *map[int32]SubType               `thrift:"lp,24,optional"`
	Mv map[int32]map[int32]string       `thrift:"mv,25,optional"`
	Mp *map[int32]map[int32]string      `thrift:"mp,26,optional"`
	Nv [][]string                       `thrift:"nv,27,optional"`
	Np *[][]string                      `thrift:"np,28,optional"`
	Ov Int                              `thrift:"ov,29"`
	Op *Int                             `thrift:"op,30,optional"`
	Pv Enum                             `thrift:"pv,31"`
	Pp *Enum                            `thrift:"pp,32,optional"`
	Qv map[int32][]string               `thrift:"qv,33,optional"`
	Qp *map[int32][]string              `thrift:"qp,34,optional"`
	Rv []map[string][]map[string]int32  `thrift:"rv,35,optional"`
	Rp *[]map[string][]map[string]int32 `thrift:"rp,36,optional"`
}
