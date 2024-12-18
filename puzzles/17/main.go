package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type registers struct {
	a int
	b int
	c int
}

type instruction struct {
	opcode  int
	operand int
}

func convertLst[T1 any, T2 any](lst []T1, convFunc func(T1) T2) []T2 {
	result := make([]T2, len(lst))
	for i, s := range lst {
		result[i] = convFunc(s)
	}
	return result
}

func getLiteral(operand int, r *registers) int {
	switch operand {
	case 4:
		return r.a
	case 5:
		return r.b
	case 6:
		return r.c
	}

	return operand
}

func getInstruction(instStr string, ptr int) instruction {
	numberStrLst := strings.Split(instStr, ",")
	opcode, _ := strconv.Atoi(numberStrLst[ptr])
	operand, _ := strconv.Atoi(numberStrLst[ptr+1])
	return instruction{opcode, operand}
}

func perform(instrStr string, r *registers) []int {
	outputs := make([]int, 0)
	instrPtr := 0

	for instrPtr < (len(strings.Split(instrStr, ",")) - 1) {
		instr := getInstruction(instrStr, instrPtr)

		switch instr.opcode {
		case 0:
			nmr := float64(r.a)
			dnm := math.Pow(2., float64(getLiteral(instr.operand, r)))
			r.a = int(math.Floor(nmr / dnm))
		case 1:
			r.b = r.b ^ int(instr.operand)
		case 2:
			r.b = getLiteral(instr.operand, r) % 8
		case 3:
			if r.a != 0 {
				instrPtr = int(instr.operand)
				continue
			}
		case 4:
			r.b = r.b ^ r.c
		case 5:
			out := getLiteral(instr.operand, r) % 8
			outputs = append(outputs, out)
		case 6:
			nmr := float64(r.a)
			dnm := math.Pow(2., float64(getLiteral(instr.operand, r)))
			r.b = int(math.Floor(nmr / dnm))
		case 7:
			nmr := float64(r.a)
			dnm := math.Pow(2., float64(getLiteral(instr.operand, r)))
			r.c = int(math.Floor(nmr / dnm))

		}

		instrPtr += 2
	}

	return outputs
}

func parseRegister(input string) int {
	_, valueStr, _ := strings.Cut(input, ": ")
	value, _ := strconv.Atoi(valueStr)
	return value
}

func parseInput(input string) (*registers, string) {
	registersStr, programStr, _ := strings.Cut(input, "\n\n")
	registersLst := strings.Split(registersStr, "\n")
	r := &registers{
		a: parseRegister(registersLst[0]),
		b: parseRegister(registersLst[1]),
		c: parseRegister(registersLst[2]),
	}

	_, programInput, _ := strings.Cut(programStr, ": ")
	return r, programInput
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

	registers, instrStr := parseInput(content)
	output := perform(instrStr, registers)
	outputStr := convertLst[int, string](
		output,
		func(i int) string {
			return strconv.Itoa(i)
		},
	)
	fmt.Println("part 1", strings.Join(outputStr, ","))

	var findA func(match int, a int) (int, bool)
	findA = func(match int, a int) (int, bool) {
		for i := 0; i < 8; i++ {
			newA := a<<3 + i
			registers.a = newA
			result := perform(instrStr, registers)

			if len(result) >= match &&
				slices.Equal(output[len(output)-match:], result[:match]) {
				if len(result) >= len(output) {
					return newA, true
				}

				if res, found := findA(match+1, newA); found {
					return res, true
				}
			}
		}

		return 0, false
	}

	a, _ := findA(0, 0)
	fmt.Println("part 2", a)
}
