package instruction

import (
	"yars/bv"
	"yars/intf"
)

type Instruction struct {
	name      string // used to indicate _inst currently under matching in debugger, or further use
	p         intf.ProcessorInterface
	pattern   bv.BitPattern
	operation func(_inst *Instruction)
	iImm      bv.BitVector
	shamt     bv.BitVector
	shamtw    bv.BitVector
	sImm      bv.BitVector
	bImm      bv.BitVector
	uImm      bv.BitVector
	jImm      bv.BitVector
	rd        bv.BitVector
	rs1       bv.BitVector
	rs2       bv.BitVector
}

func (_inst *Instruction) Match(inst bv.BitVector) bool {
	if !_inst.pattern.Match(inst) {
		return false
	}
	_inst.iImm = inst.Sub(31, 20)
	_inst.shamt = inst.Sub(24, 20) // yarc has 6 bits of shamt
	_inst.sImm = bv.Cat(inst.Sub(31, 25), inst.Sub(11, 7))
	_inst.bImm = bv.Cat(bv.Cat(bv.Cat(
		inst.Sub(31, 31),
		inst.Sub(7, 7)),
		inst.Sub(30, 25)),
		inst.Sub(11, 8))
	_inst.uImm = inst.Sub(31, 12)
	_inst.jImm = bv.Cat(bv.Cat(bv.Cat(
		inst.Sub(31, 31),
		inst.Sub(20, 20)),
		inst.Sub(19, 12)),
		inst.Sub(30, 21))
	_inst.rd = inst.Sub(11, 7)
	_inst.rs1 = inst.Sub(19, 15)
	_inst.rs2 = inst.Sub(24, 20)
	return true
}

func (_inst *Instruction) Exec(p intf.ProcessorInterface) {
	_inst.p = p
	_inst.operation(_inst)
}
