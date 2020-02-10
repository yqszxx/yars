package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// SLTIU (set less than immediate unsigned) places the value 1 in register rd if register rs1 is less
// than the sign extended immediate when both are treated as signed numbers, else 0 is written to rd.
// (The immediate is first sign-extended to 64 bits then treated as an unsigned number.)
var SLTIU = Instruction{
	name:    "SLTIU",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 011 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLTIU x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.iImm.SignExtendTo(64).ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)    // first operator of alu
		op2 := _inst.iImm.SignExtendTo(64)   // second operator of alu
		result := bv.Bv(64)                  // will hold the result of alu
		if op1.ToUint64() < op2.ToUint64() { // if unsigned_int(r1) < unsigned_int(sign_extend(iImm))
			result.Set(0) // result = 1
		} // else result = 0
		_inst.p.WriteReg(_inst.rd, result) // write result back to rd
	},
}
