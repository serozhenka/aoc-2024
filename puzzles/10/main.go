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

var steps = [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func isPosInsideMap(Map [][]int, pos [2]int) bool {
	return pos[0] >= 0 && pos[0] < len(Map) &&
		pos[1] >= 0 && pos[1] < len(Map[pos[0]])
}

func getTrails(Map [][]int, pos [2]int) map[[2]int]int {
	result := map[[2]int]int{}

	for _, step := range steps {
		nextPos := [2]int{pos[0] + step[0], pos[1] + step[1]}
		if !isPosInsideMap(Map, nextPos) {
			continue
		}

		currEl := Map[pos[0]][pos[1]]
		nextEl := Map[nextPos[0]][nextPos[1]]
		if nextEl != currEl+1 {
			continue
		}
		if nextEl == 9 {
			result[nextPos] += 1
			continue
		}

		for k, v := range getTrails(Map, nextPos) {
			result[k] += v
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
	lines := strings.Split(content, "\n")

	Map := make([][]int, len(lines))
	for i, line := range lines {
		Map[i] = make([]int, len(line))
		for j, ch := range line {
			val, err := strconv.Atoi(string(ch))
			if err != nil {
				val = -1
			}
			Map[i][j] = val
		}
	}

	result1 := 0
	result2 := 0

	for i, row := range Map {
		for j, el := range row {
			if el == 0 {
				for _, numDistinctPaths := range getTrails(Map, [2]int{i, j}) {
					result1 += 1
					result2 += numDistinctPaths
				}
			}
		}
	}

	fmt.Println("part 1", result1)
	fmt.Println("part 2", result2)
}
