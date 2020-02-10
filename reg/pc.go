package reg

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

type ProgramCounter struct {
	data bv.BitVector
}

func (pc *ProgramCounter) Init() {
	pc.data = bv.Bv(32)
}

func (pc ProgramCounter) Read() bv.BitVector {
	log.Printf("Reading pc: 0x%08X", pc.data.ToUint32())

	return pc.data
}

func (pc *ProgramCounter) Write(value bv.BitVector) {
	if value.Width != 32 {
		panic("Cannot call ProgramCounter.Write with value not being a 32 bits bv.")
	}
	log.Printf("Writing pc: 0x%08X", value.ToUint32())
	pc.data = value
}
