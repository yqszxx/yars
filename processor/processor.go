package processor

import (
	"log"
	"yars/bv"
	"yars/instruction"
	"yars/mem"
	"yars/reg"
)

type Processor struct {
	Pc  *reg.ProgramCounter        // program counter
	Reg *reg.RegisterFile          // register file
	Mem *mem.Memory                // main memory
	Is  *[]instruction.Instruction // instruction set
	Npc bv.BitVector               // next pc
}

func (p *Processor) Run() {
	for {
		pc := p.Pc.Read()         // value of current pc
		inst, _ := p.Mem.Read(pc) // fetch Instruction
		if inst.Equal(bv.B("0"), false) {
			log.Println("Encountered blank mem space, stopping...")
			break
		}
		log.Printf("Fetched inst: 0x%016X -> 0x%08X", pc.ToUint64(), inst.ToUint64())
		p.Npc = bv.Bv(64)             // new pc
		p.Npc.From(pc.ToUint64() + 4) // next pc = pc + 4
		p.Decode(inst)                // decode and execute the fetched instruction
		p.Pc.Write(p.Npc)             // update pc
	}
}

func (p *Processor) Decode(inst bv.BitVector) {
	for _, _inst := range *(p.Is) { // loop through every instruction in the instruction set (IS for short)
		if _inst.Match(inst) { // if current fetched instruction is matched one of the instructions from the IS
			_inst.Exec(p) // then execute this matched instruction
			return        // ignore rest of instructions in the instruction set
			// because only one instruction in the IS can match current fetched instruction
		}
	}
	// if we reach here, no instruction in the IS matches current fetched instruction
	// this should trigger an "Illegal Instruction" exception
	// TODO: Trigger "Illegal Instruction" exception
	log.Panicln("Illegal Instruction.")
}

func (p *Processor) WritePc(npc bv.BitVector) {
	//TODO: Check alignment
	p.Npc = npc
}

func (p *Processor) GetNpc() bv.BitVector {
	return p.Npc
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
