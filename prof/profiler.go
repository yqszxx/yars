package prof

import (
	"fmt"
	"sort"
)

type Profile struct {
	instructionStat map[string]int
}

var Pr Profile

func (pr *Profile) Init() {
	pr.instructionStat = make(map[string]int)
}

func (pr *Profile) Instruction(inst string) {
	pr.instructionStat[inst] = pr.instructionStat[inst] + 1
}

func (pr *Profile) Print() {
	s := "Instruction stat:\n"
	var instructions []string
	total := 0
	for k := range pr.instructionStat {
		instructions = append(instructions, k)
		total += pr.instructionStat[k]
	}
	sort.Slice(instructions, func(i, j int) bool {
		return pr.instructionStat[instructions[i]] > pr.instructionStat[instructions[j]] // desc
	})
	s = fmt.Sprintf("%sTotal\t%5d\n", s, total)
	for _, k := range instructions {
		s = fmt.Sprintf("%s%s\t%5d\t%6.2f%%\n",
			s,
			k,
			pr.instructionStat[k],
			float64(pr.instructionStat[k])/float64(total)*100)
	}
	print(s) // to stderr
}
