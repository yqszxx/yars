package reg

import (
	"log"
	"yars/bv"
)

type ProgramCounter uint64

func (pc ProgramCounter) ReadInt() uint64 {
	return uint64(pc)
}

func (pc *ProgramCounter) WriteInt(data uint64) {
	*pc = ProgramCounter(data)
}

func (pc ProgramCounter) Read() bv.BitVector {
	_pcInt := pc.ReadInt()

	newBv := bv.Bv(64)
	newBv.From(_pcInt)

	log.Printf("Reading pc: 0x%016X", _pcInt)

	return newBv
}

func (pc *ProgramCounter) Write(data bv.BitVector) {
	if data.Width != 64 {
		panic("Cannot write PC with a bv width is not 64.")
	}
	_data := data.ToUint64()
	log.Printf("Writing pc: 0x%016X", _data)
	pc.WriteInt(_data)
}
