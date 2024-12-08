package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

func isInsideMap(pos [2]int, mapSize [2]int) bool {
	return pos[0] >= 0 && pos[0] < mapSize[0] &&
		pos[1] >= 0 && pos[1] < mapSize[1]
}

func getAntinodes(pos1 [2]int, pos2 [2]int, mapSize [2]int, anyGridPos bool) [][2]int {
	results := [][2]int{}
	if anyGridPos {
		results = append(results, pos1, pos2)
	}

	diff := [2]int{pos2[0] - pos1[0], pos2[1] - pos1[1]}
	for i := 2; ; i++ {
		antinode := [2]int{pos1[0] + i*diff[0], pos1[1] + i*diff[1]}
		if !isInsideMap(antinode, mapSize) {
			break
		}
		results = append(results, antinode)
		if !anyGridPos {
			break
		}
	}
	return results
}

func solve(content string, anyGridPos bool) int {
	lines := strings.Split(content, "\n")
	allAntennas := map[rune][][2]int{}
	mapSize := [2]int{len(lines), len(lines[0])}

	regex := regexp.MustCompile(`[a-zA-Z0-9]{1}`)
	for i, line := range strings.Split(content, "\n") {
		for _, match := range regex.FindAllStringIndex(line, len(line)) {
			startIdx, endIdx := match[0], match[1]
			frequency := rune(line[startIdx:endIdx][0])
			allAntennas[frequency] = append(allAntennas[frequency], [2]int{i, startIdx})
		}
	}

	foundAntinodes := map[[2]int]int{}
	for _, antennas := range allAntennas {
		allCombinations := [][2][2]int{}
		for i := 0; i < len(antennas); i++ {
			for j := i + 1; j < len(antennas); j++ {
				allCombinations = append(allCombinations, [2][2]int{antennas[i], antennas[j]})
			}
		}

		for _, comb := range allCombinations {
			antinodes := make([][2]int, 0)
			antinodes = slices.Concat(
				antinodes,
				getAntinodes(comb[0], comb[1], mapSize, anyGridPos),
				getAntinodes(comb[1], comb[0], mapSize, anyGridPos),
			)

			for _, antinode := range antinodes {
				foundAntinodes[antinode] = 1
			}
		}
	}

	return len(foundAntinodes)
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

	fmt.Println("part 1", solve(content, false))
	fmt.Println("part 2", solve(content, true))
}
