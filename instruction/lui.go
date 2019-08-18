package instruction

import (
	"log"
	"yars/bv"
)

// LUI (load upper immediate) places the 20-bit U-immediate into bits 31–12 of register
// rd and places zero in the lowest 12 bits. The 32-bit result is sign-extended to 64 bits.
var LUI = Instruction{
	name:    "LUI",
	pattern: bv.P("XXXXXXXXXXXXXXXXXXXX XXXXX 0110111"),
	operation: func(_inst *Instruction) {
		op1 := bv.Cat(_inst.uImm, bv.B("0000 0000 0000"))

		log.Printf("Decoding as LUI x%d, %d",
			_inst.rd.ToUint32(),
			_inst.uImm.ToUint32())

		//uImm holds number with upper immediate already in proper position and zero in the lowest 12 bits
		_inst.p.WriteReg(_inst.rd, op1) // write the sign extended result to rd
	},
}
