package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// SLLI is a logical left shift (zeros are shifted into the lower bits). The operand to be shifted
// is in rs1, and the shift amount is encoded in the lower 6 bits of the I-immediate field.
var SLLI = Instruction{
	name:    "SLLI",
	pattern: bv.P("000000 XXXXXX XXXXX 001 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SLLI x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.shamt.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.shamt                            // second operator of alu
		result := bv.Bv(64)                           // will hold the result of alu
		result.From(op1.ToUint64() << op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)            // write result back to rd
	},
}
