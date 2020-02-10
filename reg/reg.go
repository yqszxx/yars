package reg

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

type RegisterFile struct {
	data [32]bv.BitVector // RV64I has 32 integer registers
}

func (reg *RegisterFile) Init() {
	for i := 1; i < 32; i++ {
		reg.data[i] = bv.Bv(64)
	}
}

func (reg *RegisterFile) Read(n bv.BitVector) bv.BitVector {
	_n := int(n.ToUint64())

	if _n == 0 {
		return bv.Bv(64)
	}

	value := reg.data[_n]

	log.Printf("Reading x%d: %d (0x%016X)", _n, int64(value.ToUint64()), value.ToUint64())

	return value
}

func (reg *RegisterFile) Write(n bv.BitVector, value bv.BitVector) {
	if n.Width != 5 {
		log.Panic("Cannot call RegisterFile.Write with reg_no not being a 5 bits bv.")
	}
	if value.Width != 64 {
		log.Panic("Cannot call RegisterFile.Write with value not being a 64 bits bv.")
	}

	_n := int(n.ToUint64())

	if _n == 0 {
		return
	}

	log.Printf("Writing x%d: %d (0x%016X)", _n, int64(value.ToUint64()), value.ToUint64())

	reg.data[_n] = value
}
