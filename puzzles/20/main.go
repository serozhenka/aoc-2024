package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type point [2]int

var steps = [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (p1 point) add(p2 point) point {
	return point{p1[0] + p2[0], p1[1] + p2[1]}
}

func isInsideMap(pos point, mapSize [2]int) bool {
	return pos[0] >= 0 && pos[0] < mapSize[0] &&
		pos[1] >= 0 && pos[1] < mapSize[1]
}

func abs(x int) int {
	return max(x, -x)
}

func distance(pos1, pos2 point) int {
	return abs(pos1[0]-pos2[0]) + abs(pos1[1]-pos2[1])
}

func djikstra(adjacencyMap map[point]map[point]bool, startPos point) (map[point]int, map[point]point) {
	distances := map[point]int{}
	for p := range adjacencyMap {
		if p == startPos {
			distances[p] = 0
		} else {
			distances[p] = math.MaxInt
		}
	}

	parents := map[point]point{}

	visited := map[point]bool{}
	pq := make(PriorityQueue[point], 0, len(adjacencyMap))
	heap.Init(&pq)
	heap.Push(&pq, &Item[point]{
		value:    startPos,
		priority: 0,
	})

	for len(pq) > 0 {
		item := heap.Pop(&pq).(*Item[point])
		currP := item.value

		visited[currP] = true

		for ap := range adjacencyMap[currP] {
			newDistance := distances[currP] + 1
			if newDistance < distances[ap] {
				distances[ap] = newDistance
				heap.Push(&pq, &Item[point]{
					value:    ap,
					priority: distances[ap],
				})
				parents[ap] = currP
			}
		}
	}

	return distances, parents
}

func solve(rt [][]rune, startPos point, endPos point, numMoves int, saveAtLeast int) int {
	adjacencyMap := map[point]map[point]bool{}
	for i := 0; i < len(rt); i++ {
		for j := 0; j < len(rt[0]); j++ {
			pos := point{i, j}
			adjacencyMap[pos] = map[point]bool{}

			for _, step := range steps {
				newPos := pos.add(step)
				if isInsideMap(newPos, [2]int{len(rt), len(rt[0])}) && rt[newPos[0]][newPos[1]] != '#' {
					adjacencyMap[pos][newPos] = true
				}
			}
		}
	}
	distances, parents := djikstra(adjacencyMap, startPos)
	endPosDistance := distances[endPos]

	var walk func(pos point) []point
	walk = func(pos point) []point {
		result := make([]point, 0)
		result = append(result, pos)
		if pos != startPos {
			result = append(result, walk(parents[pos])...)
		}
		return result
	}

	allShorts := 0

	for _, pos := range walk(endPos) {
		if pos == endPos {
			continue
		}

		for i := range rt {
			for j := range rt {
				if rt[i][j] == '#' {
					continue
				}

				d := distance(pos, point{i, j})
				if d > numMoves {
					continue
				}

				newDistance := distances[pos] + d + (endPosDistance - distances[point{i, j}])
				if diff := endPosDistance - newDistance; diff >= saveAtLeast {
					allShorts += 1
				}
			}
		}
	}

	return allShorts
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
	racetrack := make([][]rune, len(lines))
	var startPos, endPos point

	for i, line := range lines {
		racetrack[i] = make([]rune, len(line))
		for j, ch := range line {
			racetrack[i][j] = ch

			if ch == 'S' {
				startPos = point{i, j}
			} else if ch == 'E' {
				endPos = point{i, j}
			}
		}
	}

	fmt.Println("part 1", solve(racetrack, startPos, endPos, 2, 100))
	fmt.Println("part 2", solve(racetrack, startPos, endPos, 20, 100))
}
