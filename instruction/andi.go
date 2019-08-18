package instruction

import (
	"log"
	"yars/bv"
)

// ANDI is a logical operation that performs bitwise AND on register rs1 and the sign-extended
// 12-bit immediate and place the result in rd.
var ANDI = Instruction{
	name:    "ANDI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 111 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ANDI x%d, x%d, %d",
			_inst.rd.ToUint32(),
			_inst.rs1.ToUint32(),
			_inst.iImm.SignExtendTo(32).ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)            // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)           // second operator of alu
		result := bv.Bv(32)                          // will hold the result of alu
		result.From(op1.ToUint32() & op2.ToUint32()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)           // write result back to rd
	},
}
