package processor

import (
	"github.com/yqszxx/yars/bv"
	"github.com/yqszxx/yars/instruction"
	"github.com/yqszxx/yars/intf"
	"github.com/yqszxx/yars/mem"
	"github.com/yqszxx/yars/reg"
	"log"
)

type Processor struct {
	Pc  *reg.ProgramCounter        // program counter
	Reg *reg.RegisterFile          // register file
	Csr *reg.CsrFile               // control and status register file
	Mem *mem.Memory                // main memory
	Iss *[]instruction.Instruction // instruction set
	npc bv.BitVector               // next pc
	prv uint8                      // privilege mode
}

// Upon reset, a hart’s privilege mode is set to M. The mstatus fields MIE and MPRV are reset to
// 0. The pc is set to an implementation-defined reset vector. The mcause register is set to a value
// indicating the cause of the reset. All other hart state is undefined.
func (p *Processor) Reset() {
	p.prv = intf.PrivilegeMachine

	mstatus := bv.Bv(64)
	mstatus.Reset(reg.MstatusMIE)
	mstatus.Reset(reg.MstatusMPRV)
	p.Csr.WriteByName("mstatus", mstatus)

	pc := bv.Bv(64)
	pc.From(0x80000000) // set program entry point to 0x80000000, the same as spike
	p.Pc.Write(pc)

	// The mcause values after reset have implementation-specific interpretation, but the value 0 should
	// be returned on implementations that do not distinguish different reset conditions.
	mcause := bv.Bv(64)
	mcause.From(0)
	p.Csr.WriteByName("mcause", mcause)
}

func (p *Processor) Run() {
	for {
		pc := p.Pc.Read()         // value of current pc
		inst, _ := p.Mem.Read(pc) // fetch Instruction
		log.Printf("Fetched inst: 0x%016X -> 0x%08X", pc.ToUint64(), inst.ToUint64())
		p.npc = bv.Bv(64)             // new pc
		p.npc.From(pc.ToUint64() + 4) // next pc = pc + 4
		p.Decode(inst)                // decode and execute the fetched instruction
		p.Pc.Write(p.npc)             // update pc
	}
}

func (p *Processor) Decode(inst bv.BitVector) {
	for _, _inst := range *(p.Iss) { // loop through every instruction in the instruction set (ISS for short)
		if _inst.Match(inst) { // if current fetched instruction is matched one of the instructions from the IS
			_inst.Exec(p) // then execute this matched instruction
			return        // ignore rest of instructions in the instruction set
			// because only one instruction in the IS can match current fetched instruction
		}
	}
	// if we reach here, no instruction in the IS matches current fetched instruction
	// this should trigger an "Illegal Instruction" exception
	p.GenerateException(intf.IllegalInstruction)
}

func (p *Processor) WritePc(npc bv.BitVector) {
	if npc.Test(1) || npc.Test(0) { // npc[1:0] != 0b00, instruction address misaligned
		p.GenerateException(intf.InstructionAddressMisaligned)
	}

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

func (p *Processor) HasCsrPermission(n bv.BitVector, forWriting bool) bool {
	lowestPrivilege := n.Sub(9, 8)                  // csr[9:8] encode the lowest privilege level that can access the CSR
	rwFlag := n.Sub(11, 10)                         // 00, 01, 10: read/write; 11: read-only
	if p.prv < uint8(lowestPrivilege.ToUint64()) || // no enough privilege or
		(forWriting && rwFlag.Equal(bv.B("11"), true)) { // writing a read only CSR
		return false
	} else {
		return true
	}
}

func (p *Processor) ReadCsr(n bv.BitVector, forWriting bool) (bv.BitVector, bool) {
	if p.HasCsrPermission(n, forWriting) { // has permission
		return p.Csr.Read(n), true
	} else { // doesn't have permission
		p.GenerateException(intf.IllegalInstruction)
		return bv.Bv(64), false // dummy bv that won't be used
	}
}

func (p *Processor) WriteCsr(n bv.BitVector, data bv.BitVector) {
	p.Csr.Write(n, data)
}

func (p *Processor) GetCsrName(n bv.BitVector) string {
	return p.Csr.GetName(n)
}

func (p *Processor) GenerateException(exceptionCode int) {
	if exceptionCode == intf.IllegalInstruction { // for debug
		log.Panicln("Illegal Instruction.")
	}

	// TODO: Support interrupt

	cause := bv.Bv(64)
	cause.From(exceptionCode)

	status := p.Csr.ReadByName("mstatus")

	// TODO: Some traps have tval not 0, see priv doc 4.1.11(P55)
	tval := bv.Bv(64)

	medeleg := p.Csr.ReadByName("medeleg")
	// by default we handle the exception in M Mode, except it is delegated to S Mode
	if p.prv <= intf.PrivilegeSupervisor && // if the exception is happened in U or S Mode
		medeleg.Test(exceptionCode) { // and is delegated
		// we handle the trap in S Mode
		p.WritePc(p.Csr.ReadByName("stvec"))

		p.Csr.WriteByName("scause", cause)

		p.Csr.WriteByName("sepc", p.Pc.Read())

		p.Csr.WriteByName("stval", tval)

		if status.Test(reg.MstatusSIE) {
			status.Set(reg.MstatusSPIE)
		} else {
			status.Reset(reg.MstatusSPIE)
		}

		if p.prv == intf.PrivilegeUser {
			status.Reset(reg.MstatusSPP)
		} else { // p.prv == intf.PrivilegeSupervisor
			status.Set(reg.MstatusSPP)
		}

		status.Reset(reg.MstatusSIE)

		p.Csr.WriteByName("mstatus", status)

		p.prv = intf.PrivilegeSupervisor
	} else {
		p.WritePc(p.Csr.ReadByName("mtvec"))

		p.Csr.WriteByName("mcause", cause)

		p.Csr.WriteByName("mepc", p.Pc.Read())

		p.Csr.WriteByName("mtval", tval)

		if status.Test(reg.MstatusMIE) {
			status.Set(reg.MstatusMPIE)
		} else {
			status.Reset(reg.MstatusMPIE)
		}

		if p.prv == intf.PrivilegeUser {
			status.Reset(reg.MstatusMPP1)
			status.Reset(reg.MstatusMPP0) // MPP = 0b00
		} else if p.prv == intf.PrivilegeSupervisor {
			status.Reset(reg.MstatusMPP1)
			status.Set(reg.MstatusMPP0) // MPP = 0b01
		} else { // p.prv == intf.PrivilegeMachine
			status.Set(reg.MstatusMPP1)
			status.Set(reg.MstatusMPP0) // MPP = 0b11
		}

		status.Reset(reg.MstatusMIE)

		p.Csr.WriteByName("mstatus", status)

		p.prv = intf.PrivilegeMachine
	}

}

// When executed in U-mode, S-mode, or M-mode, an ECALL generates an environment-call-from-U-mode
// exception, environment-call-from-S-mode exception, or environment-call-from-M-mode exception,
// respectively, and performs no other operation.
func (p *Processor) EnvironmentCall() {
	switch p.prv {
	case intf.PrivilegeMachine:
		p.GenerateException(intf.EnvironmentCallFromMMode)
	case intf.PrivilegeSupervisor:
		p.GenerateException(intf.EnvironmentCallFromSMode)
	case intf.PrivilegeUser:
		p.GenerateException(intf.EnvironmentCallFromUMode)
	}
}

// When executing an xRET instruction, supposing xPP holds the value y, xIE is set to xPIE;
// the privilege mode is changed to y; xPIE is set to 1; and xPP is set to U.
func (p *Processor) TrapReturn(fromPrivilege uint8) {
	status := p.Csr.ReadByName("mstatus")

	if fromPrivilege == intf.PrivilegeSupervisor {
		if p.prv >= intf.PrivilegeSupervisor {
			p.WritePc(p.Csr.ReadByName("sepc"))

			if status.Test(reg.MstatusSPP) {
				p.prv = intf.PrivilegeSupervisor
			} else {
				p.prv = intf.PrivilegeUser
			}
			status.Reset(reg.MstatusSPP)

			if status.Test(reg.MstatusSPIE) {
				status.Set(reg.MstatusSIE)
			} else {
				status.Reset(reg.MstatusSIE)
			}

			status.Set(reg.MstatusSPIE)

			p.Csr.WriteByName("mstatus", status)
		} else {
			p.GenerateException(intf.IllegalInstruction)
		}
	} else if fromPrivilege == intf.PrivilegeMachine {
		if p.prv >= intf.PrivilegeMachine {
			p.WritePc(p.Csr.ReadByName("mepc"))

			if status.Test(reg.MstatusMPP1) && status.Test(reg.MstatusMPP0) { // MPP[1:0] == 0b11, 3, machine
				p.prv = intf.PrivilegeMachine
			} else if !status.Test(reg.MstatusMPP1) && status.Test(reg.MstatusMPP0) { // MPP[1:0] == 0b01, 1, supervisor
				p.prv = intf.PrivilegeSupervisor
			} else { // MPP[1:0] == 0b00, 0, user
				p.prv = intf.PrivilegeUser
			}
			status.Reset(reg.MstatusMPP1)
			status.Reset(reg.MstatusMPP0)

			if status.Test(reg.MstatusMPIE) {
				status.Set(reg.MstatusMIE)
			} else {
				status.Reset(reg.MstatusMIE)
			}

			status.Set(reg.MstatusMPIE)

			p.Csr.WriteByName("mstatus", status)
		} else {
			p.GenerateException(intf.IllegalInstruction)
		}
	}
}

func (p *Processor) ReadMemory(address bv.BitVector, mode uint8) (bv.BitVector, bool) {
	var data, exception bv.BitVector
	var part int

	switch mode {
	case intf.DoubleWord:
		// check for alignment
		if address.Test(2) || address.Test(1) || address.Test(0) { // if address[2:0] != 0, misaligned
			p.GenerateException(intf.LoadAddressMisaligned)
			return data, false
		} else {
			// need two memory accesses
			var dataLow, dataHigh bv.BitVector
			// first low word
			dataLow, exception = p.Mem.Read(address)
			if exception.Test(0) {
				break
			}
			// then high word
			address.Set(2)
			dataHigh, exception = p.Mem.Read(address)
			if exception.Test(0) {
				break
			}
			// the concat high and low to form the data
			data = bv.Cat(dataHigh, dataLow)
		}
	case intf.Word:
		// check for alignment
		if address.Test(1) || address.Test(0) { // if address[1:0] != 0, misaligned
			p.GenerateException(intf.LoadAddressMisaligned)
			return data, false
		} else {
			data, exception = p.Mem.Read(address)
		}
	case intf.HalfWord:
		// check for alignment
		if address.Test(0) { // if address[0] != 0, misaligned
			p.GenerateException(intf.LoadAddressMisaligned)
			return data, false
		} else {
			if address.Test(1) { // address[1] == 1, high half
				part = 1
			} else { // address[1] == 0, low half
				part = 0
			}
			address.Reset(1)
			data, exception = p.Mem.Read(address)
			if exception.Test(0) { // if there is an exception, we stop here instead of truncating the data, or the BitVector.Sub will panic.
				break
			}
			if part == 1 { // part == 1, high half
				data = data.Sub(31, 16)
			} else { // part == 0, low half
				data = data.Sub(15, 0)
			}
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
		data, exception = p.Mem.Read(address)
		if exception.Test(0) { // if there is an exception, we stop here instead of truncating the data, or the BitVector.Sub will panic.
			break
		}
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

	if exception.Test(0) {
		p.GenerateException(intf.LoadAccessFault)
		return data, false
	} else {
		return data, true
	}
}

func (p *Processor) WriteMemory(address bv.BitVector, data bv.BitVector, mode uint8) {
	var exception bv.BitVector

	switch mode {
	case intf.DoubleWord:
		// check for alignment
		if address.Test(2) || address.Test(1) || address.Test(0) { // if address[2:0] != 0, misaligned
			p.GenerateException(intf.LoadAddressMisaligned)
		} else {
			// need two memory accesses
			// first store low word
			exception = p.Mem.Write(address, bv.B("1111"), data.Sub(31, 0))
			if exception.Test(0) {
				break
			}
			// then high word
			address.Set(2)
			exception = p.Mem.Write(address, bv.B("1111"), data.Sub(63, 32))
		}
	case intf.Word:
		// check for alignment
		if address.Test(1) || address.Test(0) { // if address[1:0] != 0, misaligned
			p.GenerateException(intf.StoreAMOAddressMisaligned)
		} else {
			exception = p.Mem.Write(address, bv.B("1111"), data)
		}
	case intf.HalfWord:
		// check for alignment
		if address.Test(0) { // if address[0] != 0, misaligned
			p.GenerateException(intf.StoreAMOAddressMisaligned)
		} else {
			if address.Test(1) { // address[1] == 1, update high half
				address.Reset(1)
				exception = p.Mem.Write(address, bv.B("1100"), bv.Cat(data, bv.Bv(16))) // fill low 16 bit0 with 0
			} else { // address[1] == 0, update low half
				exception = p.Mem.Write(address, bv.B("0011"), bv.Cat(bv.Bv(16), data)) // fill high 16 bits with 0
			}
		}
	case intf.Byte: // byte is always aligned
		if address.Test(1) && address.Test(0) { // address[1:0] == 0b11, highest byte
			address.Reset(1)
			address.Reset(0)
			exception = p.Mem.Write(address, bv.B("1000"), bv.Cat(data, bv.Bv(24))) // fill low 24 bits with 0
		} else if address.Test(1) && !address.Test(0) { // address[1:0] == 0b10, second high byte
			address.Reset(1)
			exception = p.Mem.Write(address, bv.B("0100"), bv.Cat(bv.Cat(bv.Bv(8), data), bv.Bv(16))) // fill high 8 bits and low 16 bits with 0
		} else if !address.Test(1) && address.Test(0) { // address[1:0] == 0b01, second low byte
			address.Reset(0)
			exception = p.Mem.Write(address, bv.B("0010"), bv.Cat(bv.Cat(bv.Bv(16), data), bv.Bv(8))) // fill high 16 bits and low 8 bits with 0
		} else { // address[1:0] == 0b00, lowest byte
			exception = p.Mem.Write(address, bv.B("0001"), bv.Cat(bv.Bv(24), data)) // fill high 24 bits with 0
		}
	}

	if exception.Test(0) {
		p.GenerateException(intf.StoreAMOAccessFault)
	}
}
