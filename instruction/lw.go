package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// LW loads a 32-bit value from memory, then sign-extends to 64-bits before storing in rd.
var LW = Instruction{
	name:    "LW",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 010 XXXXX 0000011"),
	operation: func(_inst *Instruction) {
		// TODO: Check this
		//if _inst.rd.Equal(bv.B("00000"), true) { // load with rd==x0 is illegal
		//	_inst.p.GenerateException(intf.IllegalInstruction)
		//	return
		//}

		log.Printf("Decoding as LW x%d, %d(x%d)",
			_inst.rd.ToUint32(),
			int32(_inst.iImm.SignExtendTo(32).ToUint32()),
			_inst.rs1.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.iImm.SignExtendTo(32)            // second operator of alu
		address := bv.Bv(32)                          // will hold the result of alu
		address.From(op1.ToUint32() + op2.ToUint32()) // compute address
		result := _inst.p.ReadMemory(address, intf.Word)

		_inst.p.WriteReg(_inst.rd, result.SignExtendTo(32)) // write result back to rd
	},
}
