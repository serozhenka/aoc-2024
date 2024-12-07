package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func part1(content string) int {
	guardMap, guardPos := buildMap(content)

	direction := 0 // 0 - up, 1 - right, 2 - down, 3 - bottom
	cellsSeen := map[[2]int]int{}

	for {
		cellsSeen[guardPos] = 0
		step := getStep(direction)
		nextPos := [2]int{guardPos[0] + step[0], guardPos[1] + step[1]}
		if !isPosInsideMap(guardMap, nextPos) {
			break
		} else if guardMap[nextPos[0]][nextPos[1]] == '#' {
			direction = turnRight(direction)
			continue
		}
		guardPos = nextPos
	}
	return len(cellsSeen)
}

func part2(content string) int {
	guardMap, guardPos := buildMap(content)
	direction := 0

	doesObstacleCauseLoop := func(obstaclePos [2]int) bool {
		currPos := guardPos
		newDirection := direction
		seenCoords := map[[3]int]int{}

		for i := 0; ; i++ {
			coords := [3]int{currPos[0], currPos[1], newDirection}
			if seenCoords[coords] == 1 {
				return true
			}
			seenCoords[coords] = 1

			nextPos := sumPos(currPos, getStep(newDirection))
			if !isPosInsideMap(guardMap, nextPos) {
				return false
			} else if guardMap[nextPos[0]][nextPos[1]] == '#' || nextPos == obstaclePos {
				newDirection = turnRight(newDirection)
				continue
			}

			currPos = nextPos
		}
	}

	obstaclesChecked := map[[2]int]int{}
	result := 0

	for {
		nextPos := sumPos(guardPos, getStep(direction))

		if !isPosInsideMap(guardMap, nextPos) {
			break
		} else if guardMap[nextPos[0]][nextPos[1]] == '#' {
			direction = turnRight(direction)
			continue
		}

		if obstaclesChecked[nextPos] == 0 && doesObstacleCauseLoop(nextPos) {
			result += 1
		}
		obstaclesChecked[nextPos] = 1
		guardPos = nextPos
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
