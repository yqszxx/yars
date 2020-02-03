package intf

import (
	"yars/bv"
)

const (
	Byte     = 1
	HalfWord = 2
	Word     = 3
)

type ProcessorInterface interface {
	ReadReg(n bv.BitVector) bv.BitVector
	WriteReg(n bv.BitVector, data bv.BitVector)
	ReadPc() bv.BitVector
	GetNpc() bv.BitVector
	WritePc(npc bv.BitVector)
	ReadMemory(address bv.BitVector, mode uint8) bv.BitVector
	WriteMemory(address bv.BitVector, data bv.BitVector, mode uint8)
}
