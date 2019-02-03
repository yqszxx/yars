package mem

import (
	"fmt"
	"log"
	"os"
	"sort"
	"yars/bv"
)

type memoryCell struct {
	value uint32
}

type Memory struct {
	data map[uint64]memoryCell
}

// since this memory can only accessed with address aligned to 4 bytes,
// key represents high 62 bits of address

var fromhost, tohost uint64

func SetHtif(_fromhost uint64, _tohost uint64) {
	fromhost = _fromhost
	tohost = _tohost
}

func (mem *Memory) ReadInt(address uint64) (uint32, bool) {
	if address%4 != 0 {
		return 0, true
	}
	data := mem.data[address>>2].value
	return data, false
}

func (mem *Memory) Read(address bv.BitVector) (bv.BitVector, bv.BitVector) {
	_address := address.ToUint64()

	log.Printf("Reading memory address: 0x%016X", _address)

	data, exception := mem.ReadInt(_address)
	dataBv := bv.Bv(32)
	dataBv.From(data)
	exceptionBv := bv.Bv(1)
	exceptionBv.From(exception)

	log.Printf("with data 0x%08X and exception %v", data, exception)

	return dataBv, exceptionBv
}

func (mem *Memory) WriteInt(address uint64, mask uint8, data uint32) bool {
	if address%4 != 0 {
		return true
	}

	if address == tohost {
		if data == 1 {
			log.Printf("PASS!!!")
			os.Exit(0)
		} else {
			testCase := data >> 1
			log.Printf("FAIL!!! case:%d", testCase)
			os.Exit(-1)
		}
	}

	var _mask uint32
	var _masked uint32 = 0xFF
	var i uint = 0
	for ; i < 4; i++ {
		if (mask>>i)&1 == 1 {
			_mask |= _masked << (i * 8)
		}
	}
	mem.data[address>>2] = memoryCell{value: data | (^_mask & mem.data[address>>2].value)}
	return false
}

func (mem *Memory) Write(address bv.BitVector, mask bv.BitVector, data bv.BitVector) bv.BitVector {
	if address.Width != 64 {
		log.Panic("Cannot call Memory.Write with address not being a 64 bits bv")
	}
	if mask.Width != 4 {
		log.Panic("Cannot call Memory.Write with mask not being a 4 bits bv")
	}
	if data.Width != 32 {
		log.Panic("Cannot call Memory.Write with data not being a 32 bits bv")
	}
	_address := address.ToUint64()
	_mask := uint8(mask.ToUint64())
	_data := uint32(data.ToUint64())

	log.Printf("Writing memory: 0x%016X <- 0x%08X with mask %s", _address, _data, mask)

	exception := mem.WriteInt(_address, _mask, _data)
	exceptionBv := bv.Bv(1)
	exceptionBv.From(exception)

	log.Printf("with exception %v", exception)

	return exceptionBv
}

func (mem *Memory) Init() {
	mem.data = make(map[uint64]memoryCell)
}

func (mem Memory) String() string {
	s := "Dumping memory\n"
	var keys []uint64
	for k := range mem.data {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		s = fmt.Sprintf("%sM:0x%016X->0x%08X\n", s, k<<2, mem.data[k].value)
	}
	return s
}
