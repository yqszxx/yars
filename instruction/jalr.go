package instruction

import (
	"log"
	"yars/bv"
)

// JALR (jump and link register) uses the I-type encoding. The target address is obtained
// by adding the 12-bit signed I-immediate to the register rs1, then setting the least-significant
// bit of the result to zero. The address of the instruction following the jump (pc+4) is written
// to register rd.
var JALR = Instruction{
	name:    "JALR",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 000 XXXXX 1100111"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as JALR x%d, x%d, %d",
			_inst.rd.ToUint64(),
			_inst.rs1.ToUint64(),
			int32(_inst.iImm.SignExtendTo(64).ToUint64()))

		op1 := _inst.iImm.SignExtendTo(64) // sign extends iImm to 64 bit to form op1
		op2 := _inst.p.ReadReg(_inst.rs1)  // fetch op2 from rs1

		tmp := _inst.p.GetNpc()                   // tmp holds the address of the instruction following the jump (pc+4)
		npc := bv.Bv(64)                          // npc holds the jump target address
		npc.From(op1.ToUint64() + op2.ToUint64()) // the jump target address is formed by adding iImm to the value of rs1
		npc.Reset(0)                              // setting the least-significant bit of the result to zero.
		_inst.p.WritePc(npc)                      // set the jump target address
		_inst.p.WriteReg(_inst.rd, tmp)           // write result back to rd
	},
}
