package processor

import (
	"github.com/yqszxx/yars/bv"
	"github.com/yqszxx/yars/instruction"
	"github.com/yqszxx/yars/intf"
	"github.com/yqszxx/yars/mem"
	"github.com/yqszxx/yars/prof"
	"github.com/yqszxx/yars/reg"
	"log"
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
	var part int

	switch mode {
	case intf.Word:
		data = p.Mem.Read(address)
	case intf.HalfWord:
		if address.Test(1) { // address[1] == 1, high half
			part = 1
		} else { // address[1] == 0, low half
			part = 0
		}
		address.Reset(1)
		data = p.Mem.Read(address)
		if part == 1 { // part == 1, high half
			data = data.Sub(31, 16)
		} else { // part == 0, low half
			data = data.Sub(15, 0)
		}
	case intf.Byte: // byte is always aligned
		if address.Test(1) && address.Test(0) { // address[1:0] == 0b11, highest byte
			part = 3
		} else if address.Test(1) && !address.Test(0) { // address[1:0] == 0b10, second high byte
			part = 2
		} else if !address.Test(1) && address.Test(0) { // address[1:0] == 0b01, second low byte
			part = 1
		} else { // address[1:0] == 0b00, lowest byte
			part = 0
		}
		address.Reset(1)
		address.Reset(0)
		data = p.Mem.Read(address)
		if part == 3 { // part == 3, highest byte
			data = data.Sub(31, 24)
		} else if part == 2 { // part == 2, second high byte
			data = data.Sub(23, 16)
		} else if part == 1 { // part == 1, second low byte
			data = data.Sub(15, 8)
		} else { // part == 0, lowest byte
			data = data.Sub(7, 0)
		}
	}

	return data
}

func (p *Processor) WriteMemory(address bv.BitVector, data bv.BitVector, mode uint8) {
	switch mode {
	case intf.Word:
		p.Mem.Write(address, bv.B("1111"), data)
	case intf.HalfWord:
		if address.Test(1) { // address[1] == 1, update high half
			address.Reset(1)
			p.Mem.Write(address, bv.B("1100"), bv.Cat(data, bv.Bv(16))) // fill low 16 bit0 with 0
		} else { // address[1] == 0, update low half
			p.Mem.Write(address, bv.B("0011"), bv.Cat(bv.Bv(16), data)) // fill high 16 bits with 0
		}
	case intf.Byte: // byte is always aligned
		if address.Test(1) && address.Test(0) { // address[1:0] == 0b11, highest byte
			address.Reset(1)
			address.Reset(0)
			p.Mem.Write(address, bv.B("1000"), bv.Cat(data, bv.Bv(24))) // fill low 24 bits with 0
		} else if address.Test(1) && !address.Test(0) { // address[1:0] == 0b10, second high byte
			address.Reset(1)
			p.Mem.Write(address, bv.B("0100"), bv.Cat(bv.Cat(bv.Bv(8), data), bv.Bv(16))) // fill high 8 bits and low 16 bits with 0
		} else if !address.Test(1) && address.Test(0) { // address[1:0] == 0b01, second low byte
			address.Reset(0)
			p.Mem.Write(address, bv.B("0010"), bv.Cat(bv.Cat(bv.Bv(16), data), bv.Bv(8))) // fill high 16 bits and low 8 bits with 0
		} else { // address[1:0] == 0b00, lowest byte
			p.Mem.Write(address, bv.B("0001"), bv.Cat(bv.Bv(24), data)) // fill high 24 bits with 0
		}
	}

}
