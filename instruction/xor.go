package instruction

import (
	"log"
	"yars/bv"
)

// XOR performs bitwise XOR logical operation.
var XOR = Instruction{
	name:    "XOR",
	pattern: bv.P("0000000 XXXXX XXXXX 100 XXXXX 0110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as XOR x%d, x%d, x%d",
			_inst.rd.ToUint32(),
			_inst.rs1.ToUint32(),
			_inst.rs2.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)            // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2)            // second operator of alu
		result := bv.Bv(32)                          // will hold the result of alu
		result.From(op1.ToUint32() ^ op2.ToUint32()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)           // write result back to rd
	},
}
