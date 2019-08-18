package instruction

import (
	"log"
	"yars/bv"
)

// SLTU (set less than) performs unsigned compare, writing 1 to rd if rs1 < rs2, 0 otherwise.
var SLTU = Instruction{
	name:    "SLTU",
	pattern: bv.P("0000000 XXXXX XXXXX 011 XXXXX 0110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLTU x%d, x%d, x%d",
			_inst.rd.ToUint32(),
			_inst.rs1.ToUint32(),
			_inst.rs2.ToUint32())

		op1 := _inst.p.ReadReg(_inst.rs1)    // first operator of alu
		op2 := _inst.p.ReadReg(_inst.rs2)    // second operator of alu
		result := bv.Bv(32)                  // will hold the result of alu
		if op1.ToUint32() < op2.ToUint32() { // if unsigned_int(r1) < unsigned_int(r2)
			result.Set(0) // result = 1
		} // else result = 0
		_inst.p.WriteReg(_inst.rd, result) // write result back to rd
	},
}
