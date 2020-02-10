package instruction

import (
	"github.com/yqszxx/yars/bv"
	"github.com/yqszxx/yars/intf"
	"log"
)

// The SH instruction stores 16-bit value from the low bits of register rs2 to memory.
var SH = Instruction{
	name:    "SH",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 001 XXXXX 0100011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SH x%d, %d(x%d)",
			_inst.rs2.ToUint32(),
			int64(_inst.sImm.SignExtendTo(32).ToUint32()),
			_inst.rs1.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.sImm.SignExtendTo(32)            // second operator of alu
		address := bv.Bv(32)                          // will hold the result of alu
		address.From(op1.ToUint32() + op2.ToUint32()) // compute address
		data := _inst.p.ReadReg(_inst.rs2)
		_inst.p.WriteMemory(address, data.Sub(15, 0), intf.HalfWord)
	},
}
