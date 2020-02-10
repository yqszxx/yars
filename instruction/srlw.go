package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// SRLW is an RV64I-only instruction that is analogously defined but operates on 32-bit values
// and produces signed 32-bit result. The shift amount is given by rs2[4:0].
var SRLW = Instruction{
	name:    "SRLW",
	pattern: bv.P("0000000 XXXXX XXXXX 101 XXXXX 0111011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SRLW x%d, x%d, x%d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.rs2.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)                             // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2).Sub(4, 0)                   // second operator of alu
		result := bv.Bv(32)                                           // will hold the result of alu
		result.From(uint32(op1.ToUint64()) >> uint32(op2.ToUint64())) // perform computation
		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(64))           // write result back to rd
	},
}
