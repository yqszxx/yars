package instruction

import (
	"log"
	"yars/bv"
)

// The CSRRS (Atomic Read and Set Bits in CSR) instruction reads the value of the CSR, zero-
// extends the value to XLEN bits, and writes it to integer register rd. The initial value in integer
// register rs1 is treated as a bit mask that specifies bit positions to be set in the CSR. Any bit that
// is high in rs1 will cause the corresponding bit to be set in the CSR, if that CSR bit is writable.
// Other bits in the CSR are unaffected (though CSRs might have side effects when written).
var CSRRS = Instruction{
	name:    "CSRRS",
	pattern: bv.P("XXXXXXXXXXXX XXXXX 010 XXXXX 1110011"),
	operation: func(_inst *Instruction) {
		log.Printf("Decoding as CSRRS x%d, %s, x%d",
			_inst.rd.ToUint64(),
			_inst.p.GetCsrName(_inst.csr),
			_inst.rs1.ToUint64())

		// if rs1=x0, then the instruction will not write to the CSR at all, and
		// so shall not cause any of the side effects that might otherwise occur on a CSR write, such as raising
		// illegal instruction exceptions on accesses to read-only CSRs.
		forWriting := !_inst.rs1.Equal(bv.Bv(5), true)

		oldCsr, success := _inst.p.ReadCsr(_inst.csr, forWriting)
		if !success {
			return
		}

		if forWriting {
			newCsr := bv.Bv(64)
			newCsr.From(oldCsr.ToUint64() | _inst.p.ReadReg(_inst.rs1).ToUint64())
			_inst.p.WriteCsr(_inst.csr, newCsr)
		}

		_inst.p.WriteReg(_inst.rd, oldCsr.SignExtendTo(64))
	},
}
