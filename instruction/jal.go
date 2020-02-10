package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// JAL (jump and link) uses the J-type format, where the J-immediate encodes a signed offset
// in multiples of 2 bytes. The offset is sign-extended and added to the pc to form the jump
// target address. JAL stores the address of the instruction following the jump (pc+4) into
// register rd.
var JAL = Instruction{
	name:    "JAL",
	pattern: bv.P("XXXXXXXXXXXXXXXXXXXXXXXXX 1101111"),
	operation: func(_inst *Instruction) {
		op1 := bv.Cat(_inst.jImm, bv.B("0")).SignExtendTo(64) // multiplies 2 and sign extends to 64 bit

		sign := '+'
		targetOffset := int64(op1.ToUint64())
		if targetOffset < 0 {
			sign = '-'
			targetOffset = -targetOffset
		}
		log.Printf("Decoding as JAL x%d, pc %c 0x%X",
			_inst.rd.ToUint64(),
			sign,
			targetOffset)

		tmp := _inst.p.GetNpc()                                // tmp holds the address of the instruction following the jump (pc+4)
		npc := bv.Bv(64)                                       // npc holds the jump target address
		npc.From(op1.ToUint64() + _inst.p.ReadPc().ToUint64()) // the jump target address is formed by adding op1 to the pc
		_inst.p.WritePc(npc)                                   // set the jump target address
		_inst.p.WriteReg(_inst.rd, tmp)                        // write result back to rd
	},
}
