package instruction

import (
	"log"
	"yars/bv"
)

// SRA performs arithmetic right shift on the value in register rs1 by the shift amount held in
// register rs2. Only the low 6 bits of rs2 are considered for the shift amount.
var SRA = Instruction{
	name:    "SRA",
	pattern: bv.P("0100000 XXXXX XXXXX 101 XXXXX 0110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SRA x%d, x%d, x%d",
			_inst.rd.ToUint32(),
			_inst.rs1.ToUint32(),
			_inst.rs2.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)                            // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2).Sub(5, 0)                  // second operator of alu
		result := bv.Bv(32)                                          // will hold the result of alu
		result.From(uint32(int32(op1.ToUint32()) >> op2.ToUint32())) // perform computation
		_inst.p.WriteReg(_inst.rd, result)                           // write result back to rd
	},
}
