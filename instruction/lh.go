package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// LH loads a 16-bit value from memory, then sign-extends to 64-bits before storing in rd.
var LH = Instruction{
	name:    "LH",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 001 XXXXX 0000011"),
	operation: func(_inst *Instruction) {
		if _inst.rd.Equal(bv.B("00000"), true) { // load with rd==x0 is illegal
			_inst.p.GenerateException(intf.IllegalInstruction)
			return
		}

		log.Printf("Decoding as LH x%d, %d(x%d)",
			_inst.rd.ToUint64(),
			int64(_inst.iImm.SignExtendTo(64).ToUint64()),
			_inst.rs1.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.iImm.SignExtendTo(64)            // second operator of alu
		address := bv.Bv(64)                          // will hold the result of alu
		address.From(op1.ToUint64() + op2.ToUint64()) // compute address
		result, success := _inst.p.ReadMemory(address, intf.HalfWord)
		if !success {
			return
		}

		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(64)) // write result back to rd
	},
}
