package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// BEQ branch instruction compares two registers. It takes the branch if registers rs1 and rs2 are equal.
var BEQ = Instruction{
	name:    "BEQ",
	pattern: bv.P("XXXXXXX XXXXX XXXXX 000 XXXXX 1100011"),
	operation: func(_inst *Instruction) {
		targetOffset := bv.Cat(_inst.bImm, bv.B("0")).SignExtendTo(64) // offset is formed by bImm multiplying 2 and sign extending to 64 bits

		sign := '+'
		targetOffsetInt := int64(targetOffset.ToUint64())
		if targetOffsetInt < 0 {
			sign = '-'
			targetOffsetInt = -targetOffsetInt
		}
		log.Printf("Decoding as BEQ x%d, x%d, pc %c 0x%X",
			_inst.rs1.ToUint64(),
			_inst.rs2.ToUint64(),
			sign,
			targetOffsetInt)

		op1 := _inst.p.ReadReg(_inst.rs1)
		op2 := _inst.p.ReadReg(_inst.rs2)

		if op1.ToUint64() == op2.ToUint64() { // if rs1 == rs2, takes branch
			npc := bv.Bv(64)                                                // npc holds the jump target address
			npc.From(targetOffset.ToUint64() + _inst.p.ReadPc().ToUint64()) // the jump target address is formed by adding targetOffset to the pc
			_inst.p.WritePc(npc)                                            // set the jump target address
		}
	},
}
