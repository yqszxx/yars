package instruction

import (
	"log"
	"yars/bv"
)

// The CSRRW (Atomic Read/Write CSR) instruction atomically swaps values in the CSRs and integer
// registers. CSRRW reads the old value of the CSR, zero-extends the value to XLEN bits, then writes
// it to integer register rd. The initial value in rs1 is written to the CSR.
var CSRRW = Instruction{
	name:    "CSRRW",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 001 XXXXX 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as CSRRW x%d, %s, x%d",
			_inst.rd.ToUint64(),
			_inst.p.GetCsrName(_inst.csr),
			_inst.rs1.ToUint64())

		oldCsr, success := _inst.p.ReadCsr(_inst.csr, true)
		if !success {
			return
		}

		// A CSRRW with rs1=x0 will attempt to write zero to the destination CSR.
		_inst.p.WriteCsr(_inst.csr, _inst.p.ReadReg(_inst.rs1))
		_inst.p.WriteReg(_inst.rd, oldCsr.SignExtendTo(64))
	},
}
