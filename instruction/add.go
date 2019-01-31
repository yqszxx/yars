package instruction

import (
	"log"
	"yars/bv"
)

// ADD performs addition. Overflow is ignored and the low 64 bits of results are written to the
// destination.
var ADD = Instruction{
	name:    "ADD",
	pattern: bv.P("0000000 XXXXX XXXXX 000 XXXXX 0110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ADD x%d, x%d, x%d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.rs2.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)            // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2)            // second operator of alu
		result := bv.Bv(64)                          // will hold the result of alu
		result.From(op1.ToUint64() + op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)           // write result back to rd
	},
}
