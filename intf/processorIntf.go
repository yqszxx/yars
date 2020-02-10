package intf

import (
	"github.com/yqszxx/yars/bv"
)

const (
	PrivilegeUser       = 0
	PrivilegeSupervisor = 1
	PrivilegeMachine    = 3
)

const (
	InstructionAddressMisaligned = 0
	IllegalInstruction           = 2
	LoadAddressMisaligned        = 4
	LoadAccessFault              = 5
	StoreAMOAddressMisaligned    = 6
	StoreAMOAccessFault          = 7
	EnvironmentCallFromUMode     = 8
	EnvironmentCallFromSMode     = 9
	EnvironmentCallFromMMode     = 11
)

const (
	Byte       = 1
	HalfWord   = 2
	Word       = 3
	DoubleWord = 4
)

type ProcessorInterface interface {
	ReadReg(n bv.BitVector) bv.BitVector
	WriteReg(n bv.BitVector, data bv.BitVector)
	ReadCsr(n bv.BitVector, forWriting bool) (bv.BitVector, bool)
	WriteCsr(n bv.BitVector, data bv.BitVector)
	GetCsrName(n bv.BitVector) string
	ReadPc() bv.BitVector
	GetNpc() bv.BitVector
	WritePc(npc bv.BitVector)
	EnvironmentCall()
	TrapReturn(fromPrivilege uint8)
	ReadMemory(address bv.BitVector, mode uint8) (bv.BitVector, bool)
	WriteMemory(address bv.BitVector, data bv.BitVector, mode uint8)
	GenerateException(exceptionCode int)
}
