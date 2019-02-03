package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// The SD instruction stores 64-bit value from the register rs2 to memory.
var SD = Instruction{
	name:    "SD",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 011 XXXXX 0100011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SD x%d, %d(x%d)",
			_inst.rs2.ToUint64(),
			int64(_inst.sImm.SignExtendTo(64).ToUint64()),
			_inst.rs1.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.sImm.SignExtendTo(64)            // second operator of alu
		address := bv.Bv(64)                          // will hold the result of alu
		address.From(op1.ToUint64() + op2.ToUint64()) // compute address
		data := _inst.p.ReadReg(_inst.rs2)
		_inst.p.WriteMemory(address, data, intf.DoubleWord)
	},
}
