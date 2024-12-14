package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var regex = regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

type point [2]int

type robot struct {
	pos point
	v   point
}

func (p1 point) add(p2 point) point {
	return point{p1[0] + p2[0], p1[1] + p2[1]}
}

func (r *robot) move(spaceSize [2]int) {
	newPos := r.pos.add(r.v)

	for i := 0; i <= 1; i++ {
		if newPos[i] < 0 {
			newPos[i] += spaceSize[i]
		} else if newPos[i] >= spaceSize[i] {
			newPos[i] -= spaceSize[i]
		}
	}

	r.pos = newPos
}

func (r robot) quadrant(spaceSize [2]int) int {
	middle := [2]int{(spaceSize[0] - 1) / 2, (spaceSize[1] - 1) / 2}

	if r.pos[0] < middle[0] && r.pos[1] < middle[1] {
		return 1
	} else if r.pos[0] < middle[0] && r.pos[1] > middle[1] {
		return 2
	} else if r.pos[0] > middle[0] && r.pos[1] < middle[1] {
		return 3
	} else if r.pos[0] > middle[0] && r.pos[1] > middle[1] {
		return 4
	} else {
		return -1
	}
}

func parseInput(content string) []robot {
	lines := strings.Split(content, "\n")
	robots := make([]robot, len(lines))

	for i, line := range lines {
		match := regex.FindStringSubmatch(line)

		pos := point{}
		pos[1], _ = strconv.Atoi(match[1])
		pos[0], _ = strconv.Atoi(match[2])

		v := point{}
		v[1], _ = strconv.Atoi(match[3])
		v[0], _ = strconv.Atoi(match[4])
		robots[i] = robot{pos, v}
	}

	return robots
}

func part1(content string, spaceSize [2]int) int {
	robots := parseInput(content)

	for i := 0; i < 100; i++ {
		for j := 0; j < len(robots); j++ {
			robots[j].move(spaceSize)
		}
	}

	quadrants := map[int]int{}
	for _, robot := range robots {
		quadrants[robot.quadrant(spaceSize)] += 1
	}

	result := 1
	for quadrant, count := range quadrants {
		if quadrant != -1 {
			result *= count
		}
	}

	return result
}

func getNumNeighbors(robots []robot) int {
	positions := map[point]int{}
	for _, robot := range robots {
		positions[robot.pos] = 1
	}

	result := 0
	steps := []point{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	for _, robot := range robots {
		for _, step := range steps {
			newPos := robot.pos.add(step)
			if positions[newPos] == 1 {
				result += 1
			}
		}
	}

	return result
}

func part2(content string, spaceSize [2]int) int {
	robots := parseInput(content)
	movesToNeighbors := map[int]int{}

	for i := 0; i < spaceSize[0]*spaceSize[1]; i++ {
		for j := 0; j < len(robots); j++ {
			robots[j].move(spaceSize)
		}
		movesToNeighbors[i] = getNumNeighbors(robots)
	}

	keys := make([]int, 0, len(movesToNeighbors))
	for key := range movesToNeighbors {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return movesToNeighbors[keys[i]] > movesToNeighbors[keys[j]]
	})

	return keys[0] + 1
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)
	spaceSize := [2]int{103, 101}

	fmt.Println("part 1", part1(content, spaceSize))
	fmt.Println("part 2", part2(content, spaceSize))
}
