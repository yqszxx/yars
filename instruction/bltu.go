package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// BLTU branch instruction compares two registers. It takes the branch if rs1 is less than rs2, using
// unsigned comparison.
var BLTU = Instruction{
	name:    "BLTU",
	pattern: bv.P("XXXXXXX XXXXX XXXXX 110 XXXXX 1100011"),
	operation: func(_inst *Instruction) {
		targetOffset := bv.Cat(_inst.bImm, bv.B("0")).SignExtendTo(32) // offset is formed by bImm multiplying 2 and sign extending to 64 bits

		sign := '+'
		targetOffsetInt := int32(targetOffset.ToUint32())
		if targetOffsetInt < 0 {
			sign = '-'
			targetOffsetInt = -targetOffsetInt
		}
		log.Printf("Decoding as BLTU x%d, x%d, pc %c 0x%X",
			_inst.rs1.ToUint32(),
			_inst.rs2.ToUint32(),
			sign,
			targetOffsetInt)

		op1 := _inst.p.ReadReg(_inst.rs1)
		op2 := _inst.p.ReadReg(_inst.rs2)

		if op1.ToUint32() < op2.ToUint32() { // if unsigned_int(rs1) < unsigned_int(rs2), takes branch
			npc := bv.Bv(32)                                                // npc holds the jump target address
			npc.From(targetOffset.ToUint32() + _inst.p.ReadPc().ToUint32()) // the jump target address is formed by adding targetOffset to the pc
			_inst.p.WritePc(npc)                                            // set the jump target address
		}
	},
}
