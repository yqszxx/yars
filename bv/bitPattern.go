package bv

import (
	"fmt"
	"log"
	"strings"
)

type BitPattern struct {
	bits  []int8
	Width uint
}

// 0: LSB, n: MSB
// 0: 0, 1: 1, -1: X

func Bp(length uint) BitPattern {
	var newBp BitPattern
	newBp.Width = length
	newBp.bits = make([]int8, length)
	return newBp
}

func (bp BitPattern) Match(bv BitVector) bool {
	if len(bp.bits) != len(bv.bits) {
		return false
	}
	for i := 0; i < len(bp.bits); i++ { // because opcode resides at lower 7 bits, comparison starting from LSB can speed up the matching process
		switch bp.bits[i] {
		case 1:
			if bv.bits[i] != true {
				return false
			}
		case 0:
			if bv.bits[i] != false {
				return false
			}
		default:
			continue
		}
	}
	return true
}

func (bp BitPattern) String() string {
	s := ""
	for i := len(bp.bits) - 1; i >= 0; i-- {
		if bp.bits[i] == 1 {
			s = fmt.Sprintf("%s1", s)
		} else if bp.bits[i] == 0 {
			s = fmt.Sprintf("%s0", s)
		} else {
			s = fmt.Sprintf("%sX", s)
		}
		if i%8 == 0 {
			s = fmt.Sprintf("%s ", s)
		}
	}
	s = fmt.Sprintf("%s of len(%d)\n", s, len(bp.bits))
	return s
}

func P(s string) BitPattern {
	s = strings.Replace(s, " ", "", -1)
	newBp := Bp(uint(len(s)))
	for i := len(s) - 1; i >= 0; i-- {
		var v int8
		switch s[i] {
		case '1':
			v = 1
		case '0':
			v = 0
		case 'X':
			v = -1
		default:
			log.Panicf("Identifier '%c' shoud not appear in BitPattern.", s[i])
		}
		newBp.bits[len(s)-1-i] = v
	}
	return newBp
}
