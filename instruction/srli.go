package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// SRLI is a logical right shift (zeros are shifted into the upper bits). The operand to be shifted
// is in rs1, and the shift amount is encoded in the lower 6 bits of the I-immediate field. The
// right shift type is encoded in bit 30. A 0 in bit 30 indicates a logical shift.
var SRLI = Instruction{
	name:    "SRLI",
	pattern: bv.P("000000 XXXXXX XXXXX 101 XXXXX 0010011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SRLI x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			_inst.shamt.ToUint64())

		op1 := _inst.p.ReadReg(_inst.rs1)             // first operator of alu
		op2 := _inst.shamt                            // second operator of alu
		result := bv.Bv(64)                           // will hold the result of alu
		result.From(op1.ToUint64() >> op2.ToUint64()) // perform computation
		_inst.p.WriteReg(_inst.rd, result)            // write result back to rd
	},
}
