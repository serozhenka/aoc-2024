package main

import "strings"

func isPosInsideMap(Map [][]rune, pos [2]int) bool {
	return pos[0] >= 0 && pos[0] < len(Map) &&
		pos[1] >= 0 && pos[1] < len(Map[pos[0]])
}

func sumPos(pos1 [2]int, pos2 [2]int) [2]int {
	return [2]int{pos1[0] + pos2[0], pos1[1] + pos2[1]}
}

func buildMap(content string) ([][]rune, [2]int) {
	lines := strings.Split(content, "\n")
	guardMap := make([][]rune, len(lines))
	var guardPos [2]int

	for i, line := range lines {
		guardMap[i] = []rune(line)
		if idx := strings.Index(line, "^"); idx != -1 {
			guardPos = [2]int{i, idx}
		}
	}

	return guardMap, guardPos
}

func turnRight(direction int) int {
	if direction += 1; direction > 3 {
		direction = 0
	}
	return direction
}

func getStep(direction int) [2]int {
	step := [2]int{-1, 0}
	for i := 0; i < direction; i++ {
		step = [2]int{step[1], -step[0]}
	}
	return step
}
