package instruction

import (
	"log"
	"yars/bv"
)

// The FENCE instruction is used to order device I/O and memory accesses as viewed by other RISC-V
// harts and external devices or coprocessors. Since we have only one hart here and no other auxiliary
// device, this instruction is considered as a NOP.
var FENCE = Instruction{
	name:    "FENCE",
	pattern: bv.P("0000 XXXX XXXX 00000 000 00000 0001111"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as FENCE")

		// implement as a nop
	},
}
