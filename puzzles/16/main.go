package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

type point [2]int

type vertex struct {
	pos       point
	visited   bool
	distance  int
	direction int
}

var steps = []point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (p1 point) add(p2 point) point {
	return point{p1[0] + p2[0], p1[1] + p2[1]}
}

func (p1 point) sub(p2 point) point {
	return point{p1[0] - p2[0], p1[1] - p2[1]}
}

func getDirection(step point) int {
	// 1 - up, 2 - right, 3 - down, 4 - left
	initialStep := point{-1, 0}
	for i := 0; i < 4; i++ {
		if initialStep == step {
			return i + 1
		}
		initialStep = point{initialStep[1], -initialStep[0]}
	}
	panic("Invalid step")
}

func abs(x int) int {
	return max(x, -x)
}

func turnsRequired(direction int, newDirection int) int {
	diff := abs(direction - newDirection)
	if diff > 2 {
		return diff % 2
	}
	return diff
}

type queue[T any] []T

func (q *queue[T]) Enqueue(el T) {
	*q = append(*q, el)
}

func (q *queue[T]) Size() int {
	return len(*q)
}

func (q *queue[T]) Dequeue() (T, error) {
	var zero T
	if q.Size() == 0 {
		return zero, errors.New("queue is empty")
	}

	first := (*q)[0]
	*q = (*q)[1:]
	return first, nil
}

func dijkstra(adjecencyList [][]int) ([]int, map[int][]int) {
	visited := map[int]bool{}
	parent := map[int][]int{}

	distances := make([]int, len(adjecencyList))
	for i := range distances {
		if i != 0 {
			distances[i] = math.MaxInt
		}
	}

	for range adjecencyList {
		idx := -1
		minDistance := math.MaxInt
		for i := range adjecencyList {
			if !visited[i] && distances[i] < minDistance {
				minDistance = distances[i]
				idx = i
			}
		}

		visited[idx] = true

		for aIdx := range adjecencyList {
			if adjecencyList[idx][aIdx] != 0 && !visited[aIdx] {
				distance := distances[idx] + adjecencyList[idx][aIdx]
				if distance < distances[aIdx] {
					distances[aIdx] = distance
					parent[aIdx] = []int{idx}
				}

				if distance <= distances[aIdx] {
					if !slices.Contains(parent[aIdx], idx) {
						parent[aIdx] = append(parent[aIdx], idx)
					}
				}
			}
		}
	}

	return distances, parent
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)
	rows := strings.Split(content, "\n")

	var startPos, endPos point
	posToVertex := map[point]*vertex{}

	for i, row := range rows {
		for j, ch := range row {
			if ch == 'S' {
				startPos = point{i, j}
			} else if ch == 'E' {
				endPos = point{i, j}
			}

			if ch != '#' {
				v := &vertex{pos: point{i, j}, visited: false, distance: math.MaxInt64}
				posToVertex[v.pos] = v
			}
		}
	}

	posToVertex[startPos].distance = 0
	posToVertex[startPos].direction = 2

	allVertices := make([]*vertex, 0)
	queue := queue[*vertex]{posToVertex[startPos]}
	visited := map[[3]int]int{}

	for queue.Size() > 0 {
		v, err := queue.Dequeue()
		if err != nil {
			panic("queue is empty")
		}

		key := [3]int{v.pos[0], v.pos[1], v.direction}
		if visited[key] == 1 {
			continue
		}

		visited[key] = 1
		allVertices = append(allVertices, v)

		for _, step := range steps {
			newDirection := getDirection(step)
			avPos := v.pos.add(step)
			key = [3]int{avPos[0], avPos[1], newDirection}
			if _, exists := posToVertex[avPos]; !exists || visited[key] == 1 {
				continue
			}

			var newV *vertex
			if (v.direction - newDirection) == 0 {
				newV = &vertex{avPos, false, math.MaxInt, v.direction}
			} else {
				newV = &vertex{v.pos, false, math.MaxInt, newDirection}
			}

			queue.Enqueue(newV)
		}
	}

	adjecencyList := make([][]int, len(allVertices))
	for i := range allVertices {
		adjecencyList[i] = make([]int, len(allVertices))
	}

	for i, v1 := range allVertices {
		for j, v2 := range allVertices {
			if v1 == v2 {
				continue
			}

			if v1.pos == v2.pos {
				turns := turnsRequired(v1.direction, v2.direction)
				adjecencyList[i][j] = turns * 1000
				adjecencyList[j][i] = turns * 1000
			} else if v1.direction == v2.direction {
				diff1 := v1.pos.sub(v2.pos)
				diff2 := v2.pos.sub(v1.pos)

				if steps[v1.direction-1] == diff1 {
					adjecencyList[j][i] = 1
				} else if steps[v1.direction-1] == diff2 {
					adjecencyList[i][j] = 1
				}
			}
		}
	}

	distances, parent := dijkstra(adjecencyList)

	var traversePaths func(v int, path []int) [][]int
	traversePaths = func(v int, path []int) [][]int {
		paths := make([][]int, 0)
		path = append(path, v)

		if allVertices[v].pos == startPos {
			paths = append(paths, path)
			return paths
		}

		for _, p := range parent[v] {
			for _, newPath := range traversePaths(p, path) {
				pCopy := make([]int, len(newPath))
				copy(pCopy, newPath)
				paths = append(paths, pCopy)
			}
		}
		return paths
	}

	minDist := math.MaxInt
	for i := range allVertices {
		if allVertices[i].pos != endPos {
			continue
		}
		if distances[i] < minDist {
			minDist = distances[i]
		}
	}

	uniquePos := map[point]int{}
	for i := range allVertices {
		if allVertices[i].pos != endPos || distances[i] != minDist {
			continue
		}
		for _, path := range traversePaths(i, []int{}) {
			for _, el := range path {
				uniquePos[allVertices[el].pos] = 1
			}
		}
	}

	fmt.Println("part 1", minDist)
	fmt.Println("part 2", len(uniquePos))
}
