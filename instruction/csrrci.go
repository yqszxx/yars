package instruction

import (
	"github.com/yqszxx/yars/bv"
	"log"
)

// CSRRCI variant is similar to CSRRC, except it updates the CSR using an XLEN-bit value obtained by
// zero-extending a 5-bit unsigned immediate (uImm[4:0]) field encoded in the rs1 field instead of a
// value from an integer register.
var CSRRCI = Instruction{
	name:    "CSRRCI",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 111 XXXXX 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as CSRRCI x%d, %s, %d",
			_inst.rd.ToUint64(),
			_inst.p.GetCsrName(_inst.csr),
			_inst.rs1.ToUint64())

		// if the uImm[4:0] field is zero, then this instruction will not write to the CSR,
		// and shall not cause any of the side effects that might otherwise occur on a CSR write.
		forWriting := !_inst.rs1.Equal(bv.Bv(5), true)

		oldCsr, success := _inst.p.ReadCsr(_inst.csr, forWriting)
		if !success {
			return
		}

		if forWriting {
			newCsr := bv.Bv(64)
			newCsr.From(oldCsr.ToUint64() & (^_inst.rs1.ToUint64()))
			_inst.p.WriteCsr(_inst.csr, newCsr)
		}

		_inst.p.WriteReg(_inst.rd, oldCsr.SignExtendTo(64))
	},
}
