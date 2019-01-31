package instruction

import (
	"log"
	"yars/bv"
)

// ADDIW is an RV64I-only instruction that adds the sign-extended 12-bit immediate to register rs1
// and produces the proper sign-extension of a 32-bit result in rd. Overflows are ignored and the
// result is the low 32 bits of the result sign-extended to 64 bits.
var ADDIW = Instruction{
	name:    "ADDIW",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 0011011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ADDIW x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			int32(_inst.iImm.SignExtendTo(64).ToUint64()))

		op1 := _inst.p.ReadReg(_inst.rs1)                            // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)                           // second operator of alu
		result := bv.Bv(32)                                          // will hold the result of alu
		result.From(uint32(op1.ToUint64()) + uint32(op2.ToUint64())) // perform computation
		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(64))          // write result back to rd
	},
}
