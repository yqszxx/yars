package reg

import (
	"log"
	"yars/bv"
)

type RegisterFile struct {
	data [32]bv.BitVector // yarc has 32 integer registers
}

func (reg *RegisterFile) Init() {
	for i := 1; i < 32; i++ {
		reg.data[i] = bv.Bv(32)
	}
}

func (reg *RegisterFile) Read(n bv.BitVector) bv.BitVector {
	_n := int(n.ToUint32())

	if _n == 0 {
		return bv.Bv(32)
	}

	value := reg.data[_n]

	log.Printf("Reading x%d: %d (0x%08X)", _n, int32(value.ToUint32()), value.ToUint32())

	return value
}

func (reg *RegisterFile) Write(n bv.BitVector, value bv.BitVector) {
	if n.Width != 5 {
		log.Panic("Cannot call RegisterFile.Write with reg_no not being a 5 bits bv.")
	}
	if value.Width != 32 {
		log.Panic("Cannot call RegisterFile.Write with value not being a 32 bits bv.")
	}

	_n := int(n.ToUint32())

	if _n == 0 {
		return
	}

	log.Printf("Writing x%d: %d (0x%08X)", _n, int32(value.ToUint32()), value.ToUint32())

	reg.data[_n] = value
}
