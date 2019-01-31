package instruction

import (
	"log"
	"yars/bv"
)

// SLTI (set less than immediate) places the value 1 in register rd if register rs1 is less than the
// sign extended immediate when both are treated as signed numbers, else 0 is written to rd.
var SLTI = Instruction{
	name:    "SLTI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 010 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLTI x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.iImm.SignExtendTo(64).ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)                  // first operator of alu
		op2 := _inst.iImm.SignExtendTo(64)                 // second operator of alu
		result := bv.Bv(64)                                // will hold the result of alu
		if int64(op1.ToUint64()) < int64(op2.ToUint64()) { // if signed_int(r1) < signed_int(sign_extend(iImm))
			result.Set(0) // result = 1
		} // else result = 0
		_inst.p.WriteReg(_inst.rd, result) // write result back to rd
	},
}
