package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// The SRET instruction is used to return from traps in S-mode.
var SRET = Instruction{
	name:    "SRET",
	pattern: bv.P("0001000 00010 00000 000 00000 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SRET")

		_inst.p.TrapReturn(intf.PrivilegeSupervisor) // it's convenient to access processor status such as current privilege inside the Processor struct
	},
}
