package instruction

import (
	"log"
	"yars/bv"
)

// The ECALL instruction is used to make a request to the supporting execution environment.
var ECALL = Instruction{
	name:    "ECALL",
	pattern: bv.P("000000000000 00000 000 00000 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as ECALL")

		_inst.p.EnvironmentCall() // it's convenient to access processor status such as current privilege inside the Processor struct
	},
}
