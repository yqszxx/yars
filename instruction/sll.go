package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// SLL performs logical left shift on the value in register rs1 by the shift amount held in register
// rs2. Only the low 6 bits of rs2 are considered for the shift amount.
var SLL = Instruction{
	name:    "SLL",
	pattern: bv.P("0000000 XXXXX XXXXX 001 XXXXX 0110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLL x%d, x%d, x%d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.rs2.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2).Sub(5, 0)   // second operator of alu
		result := bv.Bv(64)                           // will hold the result of alu
		result.From(op1.ToUint64() << op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)            // write result back to rd
	},
}
