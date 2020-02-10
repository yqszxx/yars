package instruction

import (
	"github.com/yqszxx/yars/bv"
	"github.com/yqszxx/yars/intf"
	"log"
)

// The SB instruction stores 8-bit value from the low bits of register rs2 to memory.
var SB = Instruction{
	name:    "SB",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 0100011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SB x%d, %d(x%d)",
			_inst.rs2.ToUint64(),
			int64(_inst.sImm.SignExtendTo(64).ToUint64()),
			_inst.rs1.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.sImm.SignExtendTo(64)            // second operator of alu
		address := bv.Bv(64)                          // will hold the result of alu
		address.From(op1.ToUint64() + op2.ToUint64()) // compute address
		data := _inst.p.ReadReg(_inst.rs2)
		_inst.p.WriteMemory(address, data.Sub(7, 0), intf.Byte)
	},
}
