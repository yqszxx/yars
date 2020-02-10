package reg

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

var csrAddress = map[string]uint16{
	"misa":      0x301,
	"mvendorid": 0xF11,
	"marchid":   0xF12,
	"mimpid":    0xF13,
	"mhartid":   0xF14,
	"mstatus":   0x300,
	"mtvec":     0x305,
	"medeleg":   0x302,
	"mideleg":   0x303,
	"mie":       0x304,
	"mcause":    0x342,
	"mepc":      0x341,
	"stvec":     0x105,
	"sepc":      0x141,
	"scause":    0x142,
	"stval":     0x143,
	"mtval":     0x343,
	"satp":      0x180,
	"pmpaddr0":  0x3A0,
	"pmpcfg0":   0x3B0,
}

const (
	MisaMXL1    = 63
	MisaMXL0    = 62
	MstatusSIE  = 1
	MstatusMIE  = 3
	MstatusSPIE = 5
	MstatusMPIE = 7
	MstatusSPP  = 8
	MstatusMPP0 = 11
	MstatusMPP1 = 12
	MstatusMPRV = 17
	MstatusSXL0 = 34
	MstatusSXL1 = 35
	MstatusUXL0 = 32
	MstatusUXL1 = 33
)

type csrCell struct {
	name     string
	value    bv.BitVector
	writable [64]bool
}

type CsrFile struct {
	data map[uint16]*csrCell
}

func (csr *CsrFile) makeCsr(name string) *csrCell {
	csr.data[csrAddress[name]] = &csrCell{
		name:  name,
		value: bv.Bv(64),
	}
	return csr.data[csrAddress[name]]
}

func (csr *CsrFile) makeCsrWithValue(name string, value bv.BitVector) *csrCell {
	csr.data[csrAddress[name]] = &csrCell{
		name:  name,
		value: value,
	}
	return csr.data[csrAddress[name]]
}

func (cell *csrCell) setWritable(hi int, lo int) *csrCell {
	if hi < lo {
		log.Panic("Cannot call csrCell.setWritable with hi < lo.")
	}
	for i := lo; i <= hi; i++ {
		cell.writable[i] = true
	}
	return cell
}

func (cell *csrCell) setWritableBit(pos int) *csrCell {
	cell.writable[pos] = true
	return cell
}

func (csr *CsrFile) Init() {
	csr.data = make(map[uint16]*csrCell)

	misaBv := bv.Bv(64)
	misaBv.Set(MisaMXL1)   // MXL[1] = 1
	misaBv.Reset(MisaMXL0) // MXL[0] = 0, MXL[1:0] = 0b10 = 2, indicate M-XLEN = 64
	misaBv.Set('I' - 'A')  // Extension[8] = 1, indicate yars supports RV64I base ISA
	misaBv.Set('S' - 'A')  // Extension[18] = 1, indicate yars supports Supervisor Mode
	misaBv.Set('U' - 'A')  // Extension[20] = 1, indicate yars supports User Mode
	csr.makeCsrWithValue("misa", misaBv)

	csr.makeCsr("mvendorid")

	csr.makeCsr("marchid")

	csr.makeCsr("mimpid")

	csr.makeCsr("mhartid")

	mstatusBv := bv.Bv(64)
	mstatusBv.Set(MstatusSXL1)   // SXL[1] = 1
	mstatusBv.Reset(MstatusSXL0) // SXL[0] = 0, SXL[1:0] = 0b10 = 2, indicate S-XLEN = 64
	mstatusBv.Set(MstatusUXL1)   // UXL[1] = 1
	mstatusBv.Reset(MstatusUXL0) // UXL[0] = 0, UXL[1:0] = 0b10 = 2, indicate U-XLEN = 64
	csr.makeCsrWithValue("mstatus", mstatusBv).
		setWritableBit(MstatusSIE).
		setWritableBit(MstatusMIE).
		setWritableBit(MstatusSPIE).
		setWritableBit(MstatusSPP).
		setWritableBit(MstatusMPP0).
		setWritableBit(MstatusMPP1).
		setWritableBit(MstatusMPRV)

	csr.makeCsr("mtvec").setWritable(63, 2) // currently mtvec[1:0] = 0, only support direct mode, TODO: Support Interrupt, see priv doc 3.1.12(P26)

	csr.makeCsr("medeleg").setWritable(63, 0)

	csr.makeCsr("mideleg").setWritable(63, 0)

	csr.makeCsr("mie").
		setWritable(1, 0).
		setWritable(5, 3).
		setWritable(9, 7).
		setWritable(11, 11)

	csr.makeCsr("mcause").setWritable(63, 0)

	csr.makeCsr("mepc").setWritable(63, 2)

	csr.makeCsr("stvec").setWritable(63, 2)

	csr.makeCsr("sepc").setWritable(63, 2)

	csr.makeCsr("scause").setWritable(63, 0)

	csr.makeCsr("stval").setWritable(63, 0)

	csr.makeCsr("mtval").setWritable(63, 0)

	csr.makeCsr("satp") // virtual memory not implemented...

	csr.makeCsr("pmpaddr0") // physical memory attributes not implemented...

	csr.makeCsr("pmpcfg0") // physical memory protection not implemented...
}

func (csr *CsrFile) Read(n bv.BitVector) bv.BitVector {
	log.Printf("Reading csr: %s -> 0x%016X", csr.GetName(n), csr.data[uint16(n.ToUint64())].value.ToUint64())

	return csr.data[uint16(n.ToUint64())].value
}

func (csr *CsrFile) doWrite(address uint16, data *bv.BitVector, mask *[64]bool) {
	for i := 0; i < 63; i++ {
		if mask[i] {
			if data.Test(i) {
				csr.data[address].value.Set(i)
			} else {
				csr.data[address].value.Reset(i)
			}
		}
	}
}

func (csr *CsrFile) Write(address bv.BitVector, data bv.BitVector) {
	if address.Width != 12 {
		log.Panic("Cannot call CsrFile.Write with csr_no not being a 12 bits bv.")
	}
	if data.Width != 64 {
		log.Panic("Cannot call CsrFile.Write with data not being a 64 bits bv.")
	}

	log.Printf("Writing csr: %s <- 0x%016X", csr.GetName(address), data.ToUint64())

	_address := uint16(address.ToUint64())
	_mask := csr.data[_address].writable

	v := &csr.data[_address].value
	// here we check for WARL fields
	switch csr.data[_address].name {
	case "mstatus":
		if v.Test(MstatusMPP1) == true && v.Test(MstatusMPP0) == false { // we don't have a privilege numbered 0b10
			_mask[MstatusMPP1] = false // so we won't write these fields
			_mask[MstatusMPP0] = false
		}
	}

	csr.doWrite(_address, &data, &_mask)
}

func (csr *CsrFile) WriteByName(name string, data bv.BitVector) {
	addressInt, ok := csrAddress[name]
	if !ok {
		log.Panicf("No such csr: %s", name)
	}
	address := bv.Bv(12)
	address.From(addressInt)
	csr.Write(address, data)
}

func (csr *CsrFile) ReadByName(name string) bv.BitVector {
	addressInt, ok := csrAddress[name]
	if !ok {
		log.Panicf("No such csr: %s", name)
	}
	address := bv.Bv(12)
	address.From(addressInt)
	return csr.Read(address)
}

func (csr *CsrFile) GetName(n bv.BitVector) string {
	return csr.data[uint16(n.ToUint64())].name
}
