package instruction

import (
	"log"
	"yars/bv"
	"yars/intf"
)

// The MRET instruction is used to return from traps in M-mode.
var MRET = Instruction{
	name:    "MRET",
	pattern: bv.P("0011000 00010 00000 000 00000 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as MRET")

		_inst.p.TrapReturn(intf.PrivilegeMachine) // it's convenient to access processor status such as current privilege inside the Processor struct
	},
}
