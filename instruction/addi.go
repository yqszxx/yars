package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// ADDI adds the sign-extended 12-bit immediate to register rs1. Arithmetic overflow is ignored and
// the result is simply the low 64 bits of the result.
var ADDI = Instruction{
	name:    "ADDI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ADDI x%d, x%d, %d",
			_inst.rd.ToUint32(),
			_inst.rs1.ToUint32(),
			int32(_inst.iImm.SignExtendTo(32).ToUint32()))

		op1 := _inst.p.ReadReg(_inst.rs1)            // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)           // second operator of alu
		result := bv.Bv(32)                          // will hold the result of alu
		result.From(op1.ToUint32() + op2.ToUint32()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)           // write result back to rd
	},
}
