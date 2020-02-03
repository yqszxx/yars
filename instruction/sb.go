package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// The SB instruction stores 8-bit value from the low bits of register rs2 to memory.
var SB = Instruction{
	name:    "SB",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 0100011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SB x%d, %d(x%d)",
			_inst.rs2.ToUint32(),
			int64(_inst.sImm.SignExtendTo(32).ToUint32()),
			_inst.rs1.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.sImm.SignExtendTo(32)            // second operator of alu
		address := bv.Bv(32)                          // will hold the result of alu
		address.From(op1.ToUint32() + op2.ToUint32()) // compute address
		data := _inst.p.ReadReg(_inst.rs2)
		_inst.p.WriteMemory(address, data.Sub(7, 0), intf.Byte)
	},
}
