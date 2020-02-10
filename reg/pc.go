package reg

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

type ProgramCounter struct {
	data bv.BitVector
}

func (pc *ProgramCounter) Init() {
	pc.data = bv.Bv(64)
}

func (pc ProgramCounter) Read() bv.BitVector {
	log.Printf("Reading pc: 0x%016X", pc.data.ToUint64())

	return pc.data
}

func (pc *ProgramCounter) Write(value bv.BitVector) {
	if value.Width != 64 {
		panic("Cannot call ProgramCounter.Write with value not being a 64 bits bv.")
	}
	log.Printf("Writing pc: 0x%016X", value.ToUint64())
	pc.data = value
}
