package reg

import (
	"log"
	"yars/bv"
)

type registerCell struct {
	value uint64
}

type RegisterFile struct {
	data []registerCell
}

func (reg *RegisterFile) Init(n int) {
	reg.data = make([]registerCell, n)
}

func (reg RegisterFile) ReadInt(n int) uint64 {
	return reg.data[n].value
}

func (reg *RegisterFile) WriteInt(n int, data uint64) {
	if n == 0 {
		return
	}
	reg.data[n].value = data
}

func (reg RegisterFile) Read(n bv.BitVector) bv.BitVector {
	_n := int(n.ToUint64())
	_regInt := reg.ReadInt(_n)

	newBv := bv.Bv(64)
	newBv.From(_regInt)

	log.Printf("Reading x%d: %d (0x%016X)", _n, int64(_regInt), _regInt)

	return newBv
}

func (reg *RegisterFile) Write(n bv.BitVector, data bv.BitVector) {
	if n.Width != 5 {
		log.Panic("Cannot call RegisterFile.Write with reg_no not being a 5 bits bv.")
	}
	if data.Width != 64 {
		log.Panic("Cannot call RegisterFile.Write with data not being a 64 bits bv.")
	}

	_n := int(n.ToUint64())
	_data := data.ToUint64()

	log.Printf("Writing x%d: %d (0x%016X)", _n, int64(_data), _data)

	reg.WriteInt(_n, _data)
}
