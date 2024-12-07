package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func strLstToIntLst(lst []string) []int {
	result := make([]int, len(lst))
	for i, s := range lst {
		val, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Error converting string to int: %v\n", err)
		}
		result[i] = val
	}
	return result
}

func allCombinations[T any](values []T, n int) [][]T {
	if n == 0 {
		return [][]T{{}}
	}

	results := make([][]T, 0)
	for _, value := range values {
		for _, combination := range allCombinations(values, n-1) {
			results = append(results, append([]T{value}, combination...))
		}
	}
	return results
}

func add(num1 int, num2 int) int {
	return num1 + num2
}

func mul(num1 int, num2 int) int {
	return num1 * num2
}

func concat(num1 int, num2 int) int {
	num, _ := strconv.Atoi(strconv.Itoa(num1) + strconv.Itoa(num2))
	return num
}

func solve(content string, operators []func(int, int) int) int {
	calibrationResult := 0

	for _, line := range strings.Split(content, "\n") {
		testValueStr, numbersStr, _ := strings.Cut(line, ": ")
		testValue, _ := strconv.Atoi(testValueStr)
		numbers := strLstToIntLst(strings.Split(numbersStr, " "))

		opCombinations := allCombinations(operators, len(numbers)-1)
		for _, opCombination := range opCombinations {
			result := numbers[0]
			for i := 1; i < len(numbers); i++ {
				result = opCombination[i-1](result, numbers[i])
			}

			if result == testValue {
				calibrationResult += testValue
				break
			}

		}

	}

	return calibrationResult
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

	fmt.Println("part 1", solve(content, []func(int, int) int{mul, add}))
	fmt.Println("part 2", solve(content, []func(int, int) int{mul, add, concat}))
}
