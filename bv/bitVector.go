package bv

import (
	"fmt"
	"log"
	"strings"
)

type BitVector struct {
	bits  []bool
	Width uint
}

// 0: LSB, n: MSB

func Bv(width uint) BitVector {
	var newBv BitVector
	newBv.Width = width
	newBv.bits = make([]bool, width)
	return newBv
}

func (bv BitVector) ToUint32() uint32 {
	var n uint32
	for i := len(bv.bits) - 1; i >= 0; i-- {
		n <<= 1
		if bv.bits[i] {
			n |= 1
		}
	}
	return n
}

func (bv *BitVector) Set(pos int) {
	bv.bits[pos] = true
}

func (bv *BitVector) Reset(pos int) {
	bv.bits[pos] = false
}

func (bv *BitVector) Test(pos int) bool {
	return bv.bits[pos]
}

func (bv BitVector) SignExtendTo(width uint) BitVector {
	if width < uint(len(bv.bits)) {
		panic("Cannot call SignExtendTo with width < origin width")
	}
	sign := bv.bits[len(bv.bits)-1]
	newBv := Bv(width)
	for i := len(newBv.bits) - 1; i >= 0; i-- {
		if i >= len(bv.bits) {
			newBv.bits[i] = sign
		} else {
			newBv.bits[i] = bv.bits[i]
		}
	}
	return newBv
}

func (bv BitVector) Sub(hi int, lo int) BitVector {
	if hi < lo {
		panic("BV.Sub must be called with hi > lo.")
	}
	if uint(hi) > bv.Width {
		panic("BV.Sub is called with hi > Width.")
	}
	if uint(lo) > bv.Width {
		panic("BV.Sub is called with lo > Width.")
	}
	newBv := Bv(uint(hi - lo + 1))
	for i := hi - lo; i >= 0; i-- {
		newBv.bits[i] = bv.bits[lo+i]
	}
	return newBv
}

func (bv BitVector) Equal(_bv BitVector, strict bool) bool {
	if strict {
		if len(bv.bits) != len(_bv.bits) {
			return false
		}
	}
	var i int
	if len(bv.bits) < len(_bv.bits) {
		i = len(bv.bits)
	} else {
		i = len(_bv.bits)
	}
	for i--; i >= 0; i-- {
		if bv.bits[i] != _bv.bits[i] {
			return false
		}
	}
	return true
}

func (bv *BitVector) From(i interface{}) {
	var _i uint32
	switch v := i.(type) {
	case uint16:
		_i = uint32(v)
	case int:
		_i = uint32(uint(v))
	case uint32:
		_i = v
	case bool:
		bv.bits[0] = v
		return
	default:
		log.Panicf("Unable to convert %T to uint32.\n", i)
	}
	var j uint32
	for j = 0; j < uint32(len(bv.bits)); j++ {
		if (_i>>j)&1 == 1 {
			bv.bits[j] = true
		} else {
			bv.bits[j] = false
		}
	}
}

func Cat(bv1 BitVector, bv2 BitVector) BitVector {
	newBv := Bv(bv1.Width + bv2.Width)
	for i, v := range bv1.bits {
		newBv.bits[uint(i)+bv2.Width] = v
	}
	for i, v := range bv2.bits {
		newBv.bits[i] = v
	}
	return newBv
}

func (bv BitVector) String() string {
	s := ""
	for i := len(bv.bits) - 1; i >= 0; i-- {
		if bv.bits[i] {
			s = fmt.Sprintf("%s1", s)
		} else {
			s = fmt.Sprintf("%s0", s)
		}
		if i%8 == 0 {
			s = fmt.Sprintf("%s ", s)
		}
	}
	s = fmt.Sprintf("%sof width %d", s, bv.Width)
	return s
}

//func (bv BitVector) String() string {
//	var w uint
//	if bv.Width > 32 {
//		w = 16
//	} else {
//		w = 8
//	}
//	format := fmt.Sprintf("0x%%0%dX of width %%d", w)
//	return fmt.Sprintf(format, bv.ToUint32(), bv.Width)
//}

func B(s string) BitVector {
	s = strings.Replace(s, " ", "", -1)
	newBv := Bv(uint(len(s)))
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case '1':
			newBv.bits[len(s)-1-i] = true
		case '0':
			newBv.bits[len(s)-1-i] = false
		default:
			log.Panicf("Identifier '%c' shoud not appear in BitVector.", s[i])
		}
	}
	return newBv
}
