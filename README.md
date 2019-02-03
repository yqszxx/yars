YARS - Yet Another RISC-V Simulator
============


What is YARS
------------
Yars is a RISC-V ISA Simulator, implements a functional model of one
(currently) RISC-V processor. The development is currently focusing on
RV64I extension only.

Why am I writing this simulator?
------------
The golden reference, spike, has its source code complicated and
difficult to read mainly because it is written in C++ which contains
a lot of magical macros. Uhh... I know it's for performance, but it
also makes someone who wants to know how a RISC-V cpu works confused.
So, I'm teaching myself to understand the ISA by reinventing the wheel.

Why Golang?
------------
Isn't it exciting to write a big project in a new language? Okay it's
just a joke... Here are the reasons why I decided to use Golang:
1. Golang has no classes, no macros, no automatic type conversions,
   which means no magic is happening, everything you see does the most
   obvious thing as it should. It's a "Code Readers Friendly Language".
2. It compiles REALLY fast, which saves you a lot of time from waiting
   the project recompiled.
3. It generates code that runs fairly fast, enough to simulate small
   code segments like 50~100 lines of asm in a reasonable time.

About Documentation
-------------
My goal is to annotate every important line so I don't have to keep
switching between the codes and the docs. But this is not done yet...

Features Implemented
-------------
1. A functional memory.
1. An elf loader which can load program from an elf file.
1. Very limited HTIF support for running
   [riscv-tests](https://github.com/riscv/riscv-tests/),
   can resolve 'fromhost' and 'tohost' symbols. Writing a 1 to memory
   address indicated by 'tohost' will cause YARC print "PASS!!!" and
   exit with code 0, write numbers other than 1 will cause YARC exit
   with 1.
1. A processor that reads instructions and executes them.
1. All unprivileged instructions.
1. Control and Status Register function and related instructions (But
 not all CSRs defined in the spec are implemented, see CSRs Section)
1. Machine, Supervisor, and Machine Privileges
1. Exceptions support

Features To Be Implemented
-------------
1. Support interruption
1. Detailed annotations

CSRs
-------------
Accessible CSRs
_____________
Access CSRs other than these will cause a panic.
* misa
* mvendorid
* marchid
* mimpid
* mhartid
* mstatus
* mtvec
* medeleg
* mideleg
* mie
* mcause
* mepc
* stvec
* sepc
* scause
* stval
* mtval
* satp
* pmpaddr0
* pmpcfg0

Really functional CSRs
____________
* misa
* mvendorid
* marchid
* mimpid
* mhartid
* mstatus
* mtvec - only support direct mode
* medeleg
* mideleg
* mie
* mcause
* mepc
* stvec
* sepc
* scause
* stval
* mtval

References
-------------
1. [User-Level ISA Specification](https://riscv.org/specifications/)
   and 
   [Privileged ISA Specification](https://riscv.org/specifications/privileged-isa)
   of RISC-V published by RISC-V Foundation Technical Committee.
1. [Spike](https://github.com/riscv/riscv-isa-sim), the official RISC-V
   simulator.
1. [Documentation of Spike (partially)](https://github.com/poweihuang17/Documentation_Spike)
   The owner of this repo really has some good works in his other repos.
1. [FPGA開発日記](http://msyksphinz.hatenablog.com/) A blog of detailed illustrations of RISC-V in Japanese.
1. [riscv-sodor](https://github.com/ucb-bar/riscv-sodor), an educational
   processor collection written in Chisel by ucb-bar.