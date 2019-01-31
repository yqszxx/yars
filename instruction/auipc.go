package instruction

import (
	"log"
	"yars/bv"
)

// AUIPC (add upper immediate to pc) appends 12 low-order zero bits to the 20-bit U-immediate,
// sign-extends the result to 64 bits, then adds it to the pc and places the result in register rd.
var AUIPC = Instruction{
	name:    "AUIPC",
	pattern: bv.P("XXXXXXXXXXXXXXXXXXXX XXXXX 0010111"),
	operation: func(_inst *Instruction) {
		op1 := bv.Cat(_inst.uImm, bv.B("0000 0000 0000")).SignExtendTo(64) // appends 12 low-order zero bits to the 20-bit U-immediate and sign extends to 64 bit

		log.Printf("Decoding as AUIPC x%d, %d",
			_inst.rd.ToUint64(),
			_inst.uImm.ToUint64())

		result := bv.Bv(64)                                       // will hold the result
		result.From(op1.ToUint64() + _inst.p.ReadPc().ToUint64()) // adds op1 to the pc
		_inst.p.WriteReg(_inst.rd, result)                        // places the result in register rd
	},
}
