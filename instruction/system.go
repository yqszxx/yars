package instruction

import (
	"log"
	"yars/bv"
)

// The SYSTEM major opcode is used to encode all privileged instructions in the RISC-V ISA.
var SYSTEM = Instruction{
	name:    "SYSTEM",
	pattern: bv.P("XXXXXXXXXXXX XXXXX XXX XXXXX 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as SYSTEM")
	},
}
