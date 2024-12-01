package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	content, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}

	leftLocations := make([]int, 0)
	rightLocations := make([]int, 0)

	for _, line := range strings.Split(string(content), "\n") {
		splitted := strings.Split(line, strings.Repeat(" ", 3))

		first, err := strconv.Atoi(splitted[0])
		if err != nil {
			log.Fatal(err)
		}
		leftLocations = append(leftLocations, first)

		second, err := strconv.Atoi(splitted[1])
		if err != nil {
			log.Fatal(err)
		}
		rightLocations = append(rightLocations, second)
	}

	sort.Ints(leftLocations)
	sort.Ints(rightLocations)

	resultPart1 := 0
	for i := 0; i < len(leftLocations); i++ {
		resultPart1 += absInt(leftLocations[i] - rightLocations[i])
	}

	fmt.Println("part 1", resultPart1)

	resultPart2 := 0
	countMap := make(map[int]int)
	for _, el := range rightLocations {
		countMap[el] += 1
	}

	for _, el := range leftLocations {
		resultPart2 += el * countMap[el]
	}

	fmt.Println("part 2", resultPart2)
}
