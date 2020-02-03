package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// LHU loads a 16-bit value from memory, then zero-extends to 32-bits before storing in rd.
var LHU = Instruction{
	name:    "LHU",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 101 XXXXX 0000011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as LHU x%d, %d(x%d)",
			_inst.rd.ToUint32(),
			int64(_inst.iImm.SignExtendTo(32).ToUint32()),
			_inst.rs1.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)            // second operator of alu
		address := bv.Bv(32)                          // will hold the result of alu
		address.From(op1.ToUint32() + op2.ToUint32()) // compute address
		data := _inst.p.ReadMemory(address, intf.HalfWord)

		result := bv.Bv(32)
		result.From(data.ToUint32())       // zero extends data to 32 bits to form result
		_inst.p.WriteReg(_inst.rd, result) // write result back to rd
	},
}
