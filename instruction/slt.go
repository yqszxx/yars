package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// SLT (set less than) performs signed compare, writing 1 to rd if rs1 < rs2, 0 otherwise.
var SLT = Instruction{
	name:    "SLT",
	pattern: bv.P("0000000 XXXXX XXXXX 010 XXXXX 0110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLT x%d, x%d, x%d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.rs2.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)                  // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2)                  // second operator of alu
		result := bv.Bv(64)                                // will hold the result of alu
		if int64(op1.ToUint64()) < int64(op2.ToUint64()) { // if signed_int(r1) < signed_int(r2)
			result.Set(0) // result = 1
		} // else result = 0
		_inst.p.WriteReg(_inst.rd, result) // write result back to rd
	},
}
