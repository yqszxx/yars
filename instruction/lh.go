package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// LH loads a 16-bit value from memory, then sign-extends to 32-bits before storing in rd.
var LH = Instruction{
	name:    "LH",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 001 XXXXX 0000011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as LH x%d, %d(x%d)",
			_inst.rd.ToUint32(),
			int64(_inst.iImm.SignExtendTo(32).ToUint32()),
			_inst.rs1.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)            // second operator of alu
		address := bv.Bv(32)                          // will hold the result of alu
		address.From(op1.ToUint32() + op2.ToUint32()) // compute address
		result := _inst.p.ReadMemory(address, intf.HalfWord)

		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(32)) // write result back to rd
	},
}
