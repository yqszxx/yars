package processor

import (
	"log"
	"yars/bv"
	"yars/instruction"
	"yars/intf"
	"yars/mem"
	"yars/prof"
	"yars/reg"
)

type Processor struct {
	Pc  *reg.ProgramCounter        // program counter
	Reg *reg.RegisterFile          // register file
	Mem *mem.Memory                // main memory
	Iss *[]instruction.Instruction // instruction set
	npc bv.BitVector               // next pc
}

// Upon reset, a hart’s privilege mode is set to M. The mstatus fields MIE and MPRV are reset to
// 0. The pc is set to an implementation-defined reset vector. The mcause register is set to a value
// indicating the cause of the reset. All other hart state is undefined.
func (p *Processor) Reset() {
	pc := bv.Bv(32)
	pc.From(0) // set program entry point to 0, the same as yarc
	p.Pc.Write(pc)
}

func (p *Processor) Run() {
	for {
		pc := p.Pc.Read()      // value of current pc
		inst := p.Mem.Read(pc) // fetch Instruction
		log.Printf("Fetched inst: 0x%08X -> 0x%08X", pc.ToUint32(), inst.ToUint32())
		p.npc = bv.Bv(32)             // new pc
		p.npc.From(pc.ToUint32() + 4) // next pc = pc + 4
		p.Decode(inst)                // decode and execute the fetched instruction
		p.Pc.Write(p.npc)             // update pc
	}
}

func (p *Processor) Decode(inst bv.BitVector) {
	for _, _inst := range *(p.Iss) { // loop through every instruction in the instruction set (ISS for short)
		if _inst.Match(inst) { // if current fetched instruction is matched one of the instructions from the IS
			prof.Pr.Instruction(_inst.GetName())
			_inst.Exec(p) // then execute this matched instruction
			return        // ignore rest of instructions in the instruction set
			// because only one instruction in the IS can match current fetched instruction
		}
	}
	log.Fatalln("Illegal instruction.")
}

func (p *Processor) WritePc(npc bv.BitVector) {
	p.npc = npc
}

func (p *Processor) GetNpc() bv.BitVector {
	return p.npc
}

func (p *Processor) ReadPc() bv.BitVector {
	return p.Pc.Read()
}

func (p *Processor) ReadReg(n bv.BitVector) bv.BitVector {
	return p.Reg.Read(n)
}

func (p *Processor) WriteReg(n bv.BitVector, data bv.BitVector) {
	p.Reg.Write(n, data)
}

func (p *Processor) ReadMemory(address bv.BitVector, mode uint8) bv.BitVector {
	var data bv.BitVector

	switch mode {
	case intf.Word:
		data = p.Mem.Read(address)
	}
	return data
}

func (p *Processor) WriteMemory(address bv.BitVector, data bv.BitVector, mode uint8) {
	switch mode {
	case intf.Word:
		p.Mem.Write(address, bv.B("1111"), data)
	}
}
