package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// ANDI is a logical operation that performs bitwise AND on register rs1 and the sign-extended
// 12-bit immediate and place the result in rd.
var ANDI = Instruction{
	name:    "ANDI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 111 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ANDI x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.iImm.SignExtendTo(64).ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)            // first operator of alu
		op2 := _inst.iImm.SignExtendTo(64)           // second operator of alu
		result := bv.Bv(64)                          // will hold the result of alu
		result.From(op1.ToUint64() & op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)           // write result back to rd
	},
}
