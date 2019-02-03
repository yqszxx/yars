package instruction

import (
	"log"
	"yars/bv"
)

// CSRRWI variant is similar to CSRRW, except it updates the CSR using an XLEN-bit value obtained by
// zero-extending a 5-bit unsigned immediate (uImm[4:0]) field encoded in the rs1 field instead of a
// value from an integer register.
var CSRRWI = Instruction{
	name:    "CSRRWI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 101 XXXXX 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as CSRRW x%d, %s, %d",
			_inst.rd.ToUint64(),
			_inst.p.GetCsrName(_inst.csr),
			_inst.rs1.ToUint64())

		oldCsr, success := _inst.p.ReadCsr(_inst.csr, true)
		if !success {
			return
		}
		newCsr := bv.Bv(64)
		newCsr.From(_inst.rs1.ToUint64()) // sign extended uImm[4:0] to 64 bit
		_inst.p.WriteCsr(_inst.csr, newCsr)
		_inst.p.WriteReg(_inst.rd, oldCsr.SignExtendTo(64))
	},
}
