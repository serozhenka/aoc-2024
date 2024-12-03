package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

func getResult(input string, multiplicationCheckEnabled bool) int {
	regex := regexp.MustCompile(`mul\((?P<first>[0-9]+),(?P<second>[0-9]+)\)|do\(\)|don't\(\)`)

	result := 0
	multiplicationEnabled := true

	for _, match := range regex.FindAllStringSubmatch(input, -1) {
		op := match[0]
		if op == "don't()" {
			multiplicationEnabled = false
			continue
		} else if op == "do()" {
			multiplicationEnabled = true
			continue
		}

		first, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}

		second, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}

		if (multiplicationCheckEnabled && multiplicationEnabled) ||
			!multiplicationCheckEnabled {
			result += first * second
		}
	}

	return result
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

	fmt.Println("part 1", getResult(content, false))
	fmt.Println("part 2", getResult(content, true))
}
