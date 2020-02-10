package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// The FENCE.I instruction is used to synchronize the instruction and data streams.
// Since we have only one instruction and data stream, this instruction is considered as a NOP.
var FENCE_I = Instruction{
	name:    "FENCE.I",
	pattern: bv.P("0000 0000 0000 00000 001 00000 0001111"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as FENCE.I")

		// implement as a nop
	},
}
