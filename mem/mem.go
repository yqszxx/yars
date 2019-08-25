package mem

import (
	"fmt"
	"log"
	"os"
	"sort"
	"yars/bv"
	"yars/prof"
)

type memoryCell struct {
	value uint32
}

type Memory struct {
	data map[uint32]memoryCell
}

// since this memory can only accessed with address aligned to 4 bytes,
// key represents high 62 bits of address

func (mem *Memory) ReadInt(address uint32) uint32 {
	data := mem.data[address>>2].value
	return data
}

func (mem *Memory) Read(address bv.BitVector) bv.BitVector {
	_address := address.ToUint32()

	log.Printf("Reading memory address: 0x%08X", _address)

	data := mem.ReadInt(_address)
	dataBv := bv.Bv(32)
	dataBv.From(data)

	log.Printf("with data 0x%08X", data)

	return dataBv
}

func (mem *Memory) WriteInt(address uint32, mask uint8, data uint32) {

	// Print Logic
	if address == 0xFFF8 {
		fmt.Printf("$0x%08X$\n", data)
		return
	}

	// Done Logic
	if address == 0xFFFC {
		fmt.Printf("!!!!!!!!!!!!!!!!!!!DONE#0x%08X#!!!!!!!!!!!!!!!!!!!\n", data)
		prof.Pr.Print()
		os.Exit(0)
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
}

func (mem *Memory) Write(address bv.BitVector, mask bv.BitVector, data bv.BitVector) {
	if address.Width != 32 {
		log.Panic("Cannot call Memory.Write with address not being a 32 bits bv")
	}
	if mask.Width != 4 {
		log.Panic("Cannot call Memory.Write with mask not being a 4 bits bv")
	}
	if data.Width != 32 {
		log.Panic("Cannot call Memory.Write with data not being a 32 bits bv")
	}
	_address := address.ToUint32()
	_mask := uint8(mask.ToUint32())
	_data := uint32(data.ToUint32())

	mem.WriteInt(_address, _mask, _data)

	log.Printf("Writing memory: 0x%08X <- 0x%08X with mask %s", _address, _data, mask)
}

func (mem *Memory) Init() {
	mem.data = make(map[uint32]memoryCell)
}

func (mem Memory) String() string {
	s := "Dumping memory\n"
	var keys []uint32
	for k := range mem.data {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		s = fmt.Sprintf("%sM:0x%08X -> 0x%08X\n", s, k<<2, mem.data[k].value)
	}
	s = fmt.Sprintf("%s+========\n", s)
	for _, k := range keys {
		s = fmt.Sprintf("%s%08X\n", s, mem.data[k].value)
	}
	s = fmt.Sprintf("%s-========\n", s)
	return s
}
