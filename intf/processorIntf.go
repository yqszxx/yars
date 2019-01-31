package intf

import (
	"yars/bv"
)

type ProcessorInterface interface {
	ReadReg(n bv.BitVector) bv.BitVector
	WriteReg(n bv.BitVector, data bv.BitVector)
	ReadPc() bv.BitVector
	GetNpc() bv.BitVector
	WritePc(npc bv.BitVector)
}
