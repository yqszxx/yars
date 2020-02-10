package instruction

import (
	"github.com/yqszxx/yars/bv"
	"github.com/yqszxx/yars/intf"
	"log"
)

// LB loads a 8-bit value from memory, then sign-extends to 32-bits before storing in rd.
var LB = Instruction{
	name:    "LB",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 0000011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as LB x%d, %d(x%d)",
			_inst.rd.ToUint32(),
			int64(_inst.iImm.SignExtendTo(32).ToUint32()),
			_inst.rs1.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)            // second operator of alu
		address := bv.Bv(32)                          // will hold the result of alu
		address.From(op1.ToUint32() + op2.ToUint32()) // compute address
		result := _inst.p.ReadMemory(address, intf.Byte)

		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(32)) // write result back to rd
	},
}
