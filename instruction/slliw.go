package instruction

import (
	"log"
	"yars/bv"
)

// SLLIW is an RV64I-only instruction that is analogously defined but operates on 32-bit values and
// produces signed 32-bit results. It generates an illegal instruction exception if imm[5] != 0
var SLLIW = Instruction{
	name:    "SLLIW",
	pattern: bv.P("0000000 XXXXX XXXXX 001 XXXXX 0011011"), // it won't match if imm[5] != 0 instead of generating an exception
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLLIW x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			int32(_inst.shamtw.SignExtendTo(64).ToUint64()))

		op1 := _inst.p.ReadReg(_inst.rs1)                     // first operator of alu
		op2 := _inst.shamtw                                   // second operator of alu
		result := bv.Bv(32)                                   // will hold the result of alu
		result.From(uint32(op1.ToUint64()) << op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(64))   // write result back to rd
	},
}
