package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
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

func part1(content string) int {
	splittedContent := strings.Split(content, "\n\n")

	goesAfter := map[int][]int{}
	for _, line := range strings.Split(splittedContent[0], "\n") {
		xStr, yStr, _ := strings.Cut(line, "|")
		x, _ := strconv.Atoi(xStr)
		y, _ := strconv.Atoi(yStr)
		goesAfter[x] = append(goesAfter[x], y)
	}

	result := 0

	for _, line := range strings.Split(splittedContent[1], "\n") {
		update := strLstToIntLst(strings.Split(line, ","))
		for i := 0; i < len(update); i++ {
			orderingRespected := true
			for _, mustBeAfter := range goesAfter[update[i]] {
				if slices.Contains(update[:i], mustBeAfter) {
					orderingRespected = false
					break
				}
			}

			if !orderingRespected {
				break
			}

			if i == len(update)-1 {
				result += update[len(update)/2]
			}
		}
	}

	return result
}

func part2(content string) int {
	splittedContent := strings.Split(content, "\n\n")

	goesAfter := map[int][]int{}
	goesBefore := map[int][]int{}
	for _, line := range strings.Split(splittedContent[0], "\n") {
		xStr, yStr, _ := strings.Cut(line, "|")
		x, _ := strconv.Atoi(xStr)
		y, _ := strconv.Atoi(yStr)
		goesAfter[x] = append(goesAfter[x], y)
		goesBefore[y] = append(goesBefore[y], x)
	}

	result := 0
	incorrectUpdates := make([][]int, 0)

	for _, line := range strings.Split(splittedContent[1], "\n") {
		update := strLstToIntLst(strings.Split(line, ","))

		for i := 0; i < len(update); i++ {
			orderingRespected := true
			for _, mustBeAfter := range goesAfter[update[i]] {
				if slices.Contains(update[:i], mustBeAfter) {
					orderingRespected = false
					break
				}
			}

			if !orderingRespected {
				incorrectUpdates = append(incorrectUpdates, update)
				break
			}
		}
	}

	moveAfter := func(lst []int, i int, j int) {
		// move element at index i after element at index j
		for ; i < j; i++ {
			lst[i], lst[i+1] = lst[i+1], lst[i]
		}
	}

	for _, update := range incorrectUpdates {
		for i := 0; i < len(update); {
			moved := false
			for j := i + 1; j < len(update); j++ {
				if slices.Contains(goesBefore[update[i]], update[j]) {
					moveAfter(update, i, j)
					moved = true
					break
				}
			}

			if !moved {
				i++
			}
		}
	}

	for _, update := range incorrectUpdates {
		result += update[len(update)/2]
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

	fmt.Println("part 1", part1(content))
	fmt.Println("part 2", part2(content))
}
