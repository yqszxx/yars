package main

import (
	"debug/elf"
	"fmt"
	"github.com/yqszxx/yars/instruction"
	"github.com/yqszxx/yars/mem"
	"github.com/yqszxx/yars/processor"
	"github.com/yqszxx/yars/reg"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile) // only show file name and line number

	if len(os.Args) < 2 {
		fmt.Println("Usage: yars elf_file")
		os.Exit(1)
	}

	var ram mem.Memory                            // the main memory
	ram.Init()                                    // init main memory
	fromhost, tohost := loadElf(os.Args[1], &ram) // load executable to main memory
	log.Print(ram)                                // log the content of memory
	mem.SetHtif(fromhost, tohost)

	var regs reg.RegisterFile
	regs.Init()

	var csrs reg.CsrFile
	csrs.Init()

	iss := instruction.InstructionSet // the instruction set

	var pc reg.ProgramCounter // the program counter register
	pc.Init()

	core := processor.Processor{
		Pc:  &pc,
		Reg: &regs,
		Csr: &csrs,
		Mem: &ram,
		Iss: &iss,
	} // generate one processing core

	core.Reset() // reset on power on

	core.Run() // fire the core
}

// read all loadable segment to main memory from an elf file
func loadElf(filename string, mem *mem.Memory) (uint64, uint64) {
	_elf, err := elf.Open(filename)
	if err != nil { // if we have problem opening the elf file...
		log.Panicf("Error opening executable file '%s', reason: %s", filename, err)
	}
	for _, p := range _elf.Progs { // go through every segment in elf
		if p.Type == elf.PT_LOAD && p.Memsz != 0 { // if this segment has PT_LOAD flag and will occupy space in memory
			if p.Filesz != 0 { // ...and this segment really has some stuff to load
				var i uint64                   // the offset in this segment
				for i = 0; i < p.Filesz; i++ { // loop through every bytes in this segment
					address := p.Paddr + i               // generate physic address to load this byte, paddr in the header indicate the base address to load this segment
					var buffer []byte                    // the buffer which will be loaded with data
					buffer = make([]byte, 1)             // make this buffer 1 byte in size, because we want to load the segment byte by byte
					cnt, _ := p.ReadAt(buffer, int64(i)) // read out 1 byte
					if cnt == 0 {                        // if we cannot read even 1 byte, normally this will not be true, but in case the elf is broken
						break // stop reading
					}
					data := buffer[0] // fetch the data out of the slice (this var is redundant only for making code more readable)
					var mask uint8    // start to generate the "memory updating mask", see more in mem/mem.go
					switch address & 3 {
					case 0:
						mask = 0x01
					case 1:
						mask = 0x02
					case 2:
						mask = 0x04
					case 3:
						mask = 0x08
					}
					mem.WriteInt(address & ^uint64(3), mask, uint32(data)<<((address&3)*8)) // write our 1 byte of data into main memory
				}
			}
		}
	}
	var fromhost, tohost uint64
	symbols, _ := _elf.Symbols()
	for _, v := range symbols {
		if v.Name == "fromhost" {
			fromhost = v.Value
		}
		if v.Name == "tohost" {
			tohost = v.Value
		}
	}
	return fromhost, tohost
}
