package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// ADDW is an RV64I-only instruction that is defined analogously to ADD, but operates on 32-bit value
// and produce signed 32-bit result. Overflow is ignored, and the low 32-bits of the result is
// sign-extended to 64-bits and written to the destination register.
var ADDW = Instruction{
	name:    "ADDW",
	pattern: bv.P("0000000 XXXXX XXXXX 000 XXXXX 0111011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ADDW x%d, x%d, x%d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.rs2.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)                            // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2)                            // second operator of alu
		result := bv.Bv(32)                                          // will hold the result of alu
		result.From(uint32(op1.ToUint64()) + uint32(op2.ToUint64())) // perform computation
		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(64))          // write result back to rd
	},
}
