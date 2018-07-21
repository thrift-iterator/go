package test

import (
	"github.com/stretchr/testify/require"
	"github.com/thrift-iterator/go/test"
	"github.com/thrift-iterator/go/test/level_2/struct_complex_test"
	"testing"
)

func Test_marshal_struct_complex(t *testing.T) {
	should := require.New(t)
	for _, c := range test.MarshalCombinations[:] {
		var obj struct_complex_test.TestObject
		obj.Av = false
		obj.Ap = &obj.Av
		obj.Bv = 1
		obj.Bp = &obj.Bv
		obj.Cv = 2
		obj.Cp = &obj.Cv
		obj.Dv = 3
		obj.Dp = &obj.Dv
		obj.Ev = 4
		obj.Ep = &obj.Ev
		obj.Fv = 5
		obj.Fp = &obj.Fv
		obj.Gv = 3.1415926
		obj.Gp = &obj.Gv
		obj.Hv = "6"                        // 15
		obj.Hp = &obj.Hv                    // 16
		obj.Iv = []byte{7}                  // 17
		obj.Ip = &obj.Iv                    // 18
		obj.Jv = []string{"8"}              // 19
		obj.Jp = &obj.Jv                    // 20
		obj.Kv = map[string]bool{"9": true} // 21
		obj.Kp = &obj.Kv
		obj.Lv = map[int32]struct_complex_test.SubType{10: {A: 10}}
		obj.Lp = &obj.Lv
		obj.Mv = map[int32]map[int32]string{
			101: {102: "103"},
		}
		obj.Mp = &obj.Mv
		obj.Nv = [][]string{
			{"201", "202"},
		}
		obj.Np = &obj.Nv
		obj.Ov = 11
		obj.Op = &obj.Ov
		obj.Pv = struct_complex_test.Enum_B
		obj.Pp = &obj.Pv
		obj.Qv = map[int32][]string{
			12: {"1201", "1201"},
		}
		obj.Qp = &obj.Qv
		obj.Rv = []map[string][]map[string]int32{
			{"foo": []map[string]int32{
				{"foo1": 1801},
				{"foo2": 1802},
			}},
			{"bar": []map[string]int32{
				{"bar1": 1803},
				{"bar2": 1804},
			}},
		}
		obj.Rp = &obj.Rv

		output, err := c.Marshal(obj)
		should.NoError(err)
		output1, err := c.Marshal(&obj)
		should.NoError(err)
		should.Equal(output, output1)

		var val *struct_complex_test.TestObject
		should.NoError(c.Unmarshal(output, &val))

		should.Equal(obj.Av, val.Av)
		should.Equal(*obj.Ap, *val.Ap)
		should.Equal(obj.Bv, val.Bv)
		should.Equal(obj.Bv, *val.Bp)
		should.Equal(obj.Cv, val.Cv)
		should.Equal(obj.Cv, *val.Cp)
		should.Equal(obj.Dv, val.Dv)
		should.Equal(obj.Dv, *val.Dp)
		should.Equal(obj.Ev, val.Ev)
		should.Equal(obj.Ev, *val.Ep)
		should.Equal(obj.Fv, val.Fv)
		should.Equal(obj.Fv, *val.Fp)
		should.Equal(obj.Gv, val.Gv)
		should.Equal(obj.Gv, *val.Gp)
		should.Equal(obj.Hv, val.Hv)
		should.Equal(obj.Hv, *val.Hp)
		should.Equal(obj.Iv, val.Iv)
		should.Equal(obj.Iv, *val.Ip)
		should.Equal(obj.Jv, val.Jv)
		should.Equal(obj.Jv, *val.Jp)
		should.Equal(obj.Kv, val.Kv)
		should.Equal(obj.Kv, *val.Kp)
		should.Equal(obj.Lv, val.Lv)
		should.Equal(obj.Lv, *val.Lp)
		should.Equal(obj.Mv, val.Mv)
		should.Equal(obj.Mv, *val.Mp)
		should.Equal(obj.Nv, val.Nv)
		should.Equal(obj.Nv, *val.Np)
		should.Equal(obj.Ov, val.Ov)
		should.Equal(obj.Ov, *val.Op)
		should.Equal(obj.Pv, val.Pv)
		should.Equal(obj.Pv, *val.Pp)
		should.Equal(obj.Qv, val.Qv)
		should.Equal(obj.Qv, *val.Qp)
		should.Equal(obj.Rv, val.Rv)
		should.Equal(obj.Rv, *val.Rp)
	}
}
