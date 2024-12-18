package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type point [2]int

var steps = []point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

var spaceSize = [2]int{71, 71}
var fallingBytes = 1024

func (p1 point) add(p2 point) point {
	return point{p1[0] + p2[0], p1[1] + p2[1]}
}

func isInsideMap(pos point, mapSize [2]int) bool {
	return pos[0] >= 0 && pos[0] < mapSize[0] &&
		pos[1] >= 0 && pos[1] < mapSize[1]
}

func djikstra(adjacencyMap map[point]map[point]bool) map[point]int {
	distances := map[point]int{}
	for p := range adjacencyMap {
		if p == [2]int{0, 0} {
			distances[p] = 0
		} else {
			distances[p] = math.MaxInt
		}
	}

	visited := map[point]bool{}

	for range adjacencyMap {
		var currP point
		minDistance := math.MaxInt
		for p := range adjacencyMap {
			if !visited[p] && distances[p] < minDistance {
				minDistance = distances[p]
				currP = p
			}
		}

		visited[currP] = true

		for ap := range adjacencyMap[currP] {
			if !visited[ap] {
				newDistance := distances[currP] + 1
				if newDistance < distances[ap] {
					distances[ap] = newDistance
				}
			}
		}
	}

	return distances

}

func copy2d[T any](matrix [][]T) [][]T {
	duplicate := make([][]T, len(matrix))
	for i := range matrix {
		duplicate[i] = make([]T, len(matrix[i]))
		copy(duplicate[i], matrix[i])
	}

	return duplicate
}

func sprinkleBytes(memorySpace [][]rune, bytes []point) [][]rune {
	spaceCopy := copy2d(memorySpace)
	for _, b := range bytes {
		spaceCopy[b[1]][b[0]] = '#'
	}

	return spaceCopy
}

func solve(memorySpace [][]rune) int {
	adjacencyMap := map[point]map[point]bool{}
	for i := 0; i < spaceSize[0]; i++ {
		for j := 0; j < spaceSize[1]; j++ {
			pos := point{i, j}
			adjacencyMap[pos] = map[point]bool{}

			for _, step := range steps {
				newPos := pos.add(step)
				if isInsideMap(newPos, spaceSize) && memorySpace[newPos[0]][newPos[1]] != '#' {
					adjacencyMap[pos][newPos] = true
				}
			}
		}
	}

	distances := djikstra(adjacencyMap)
	return distances[[2]int{spaceSize[0] - 1, spaceSize[1] - 1}]
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
	bytes := make([]point, len(lines))
	for i, line := range lines {
		xStr, yStr, _ := strings.Cut(line, ",")
		x, _ := strconv.Atoi(xStr)
		y, _ := strconv.Atoi(yStr)
		bytes[i] = point{x, y}
	}

	memorySpace := make([][]rune, spaceSize[0])
	for i := 0; i < spaceSize[0]; i++ {
		memorySpace[i] = make([]rune, spaceSize[1])
		for j := 0; j < spaceSize[1]; j++ {
			memorySpace[i][j] = '.'
		}
	}

	fmt.Println("part 1", solve(sprinkleBytes(memorySpace, bytes[:fallingBytes])))

	left, right := 0, len(bytes)-1
	var mid int
	forward := true

	for {
		mid = left + (right-left)/2
		fd := solve(sprinkleBytes(memorySpace, bytes[:mid]))
		var adjacentIdx int

		if forward {
			adjacentIdx = mid + 1
		} else {
			adjacentIdx = mid - 1
		}
		sd := solve(sprinkleBytes(memorySpace, bytes[:adjacentIdx]))

		if max(fd, sd) == math.MaxInt && min(fd, sd) != math.MaxInt {
			break
		} else if max(fd, sd) != math.MaxInt {
			left = mid + 1
			forward = true
		} else {
			right = mid - 1
			forward = false
		}
	}

	fmt.Printf("part 2 %d,%d\n", bytes[mid][0], bytes[mid][1])
}
