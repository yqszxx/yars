package instruction

import (
	"log"
	"yars/bv"
)

// ADDI adds the sign-extended 12-bit immediate to register rs1. Arithmetic overflow is ignored and
// the result is simply the low 64 bits of the result.
var ADDI = Instruction{
	name:    "ADDI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ADDI x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			int64(_inst.iImm.SignExtendTo(64).ToUint64()))

		op1 := _inst.p.ReadReg(_inst.rs1)            // first operator of alu
		op2 := _inst.iImm.SignExtendTo(64)           // second operator of alu
		result := bv.Bv(64)                          // will hold the result of alu
		result.From(op1.ToUint64() + op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)           // write result back to rd
	},
}
