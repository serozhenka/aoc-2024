package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type warehouse [][]rune
type point [2]int

func (p1 point) add(p2 point) point {
	return point{p1[0] + p2[0], p1[1] + p2[1]}
}

var movesMap = map[rune]point{
	'^': {-1, 0},
	'>': {0, 1},
	'v': {1, 0},
	'<': {0, -1},
}

func (w warehouse) ToString() string {
	lines := make([]string, len(w))
	for i, row := range w {
		lines[i] = string(row)
	}
	return strings.Join(lines, "\n")
}

func (w warehouse) robotPos() point {
	for i, row := range w {
		for j, item := range row {
			if item == '@' {
				return [2]int{i, j}
			}
		}
	}
	panic("unable to find robot")
}

func (w warehouse) move(moves []rune) {
	robotPos := w.robotPos()

	var pushAround func(point, point) bool
	pushAround = func(pos point, step point) bool {
		newPos := pos.add(step)
		if w.at(newPos) == 'O' {
			pushed := pushAround(newPos, step)
			if pushed {
				swap(w, newPos, pos)
			}
			return pushed
		} else if w.at(newPos) == '.' {
			swap(w, newPos, pos)
			return true
		} else {
			return false
		}
	}

	for _, move := range moves {
		step := movesMap[move]
		newPos := robotPos.add(step)

		if w.at(newPos) == '.' {
			swap(w, newPos, robotPos)
			robotPos = newPos
		} else if w.at(newPos) == '#' {
			continue
		} else {
			if pushAround(robotPos, step) {
				robotPos = newPos
			}
		}
	}
}

func (w warehouse) sumGps(symbol rune) int {
	result := 0
	for i, row := range w {
		for j, ch := range row {
			if ch == symbol {
				result += 100*i + j
			}
		}
	}

	return result
}

func (w warehouse) at(p point) rune {
	return w[p[0]][p[1]]
}

func swap[T any](lst [][]T, pos1 point, pos2 point) {
	lst[pos1[0]][pos1[1]], lst[pos2[0]][pos2[1]] = lst[pos2[0]][pos2[1]], lst[pos1[0]][pos1[1]]
}

func parseInput(content string) (warehouse, []rune) {
	warehouseStr, movesStr, _ := strings.Cut(content, "\n\n")
	warehouseLines := strings.Split(warehouseStr, "\n")

	w := make(warehouse, len(warehouseLines))
	for i, line := range warehouseLines {
		w[i] = make([]rune, len(line))
		for j, ch := range line {
			w[i][j] = ch
		}
	}

	moves := make([]rune, 0)
	for _, line := range strings.Split(movesStr, "\n") {
		for _, ch := range line {
			moves = append(moves, ch)
		}
	}

	return w, moves
}

func part1(content string) int {
	w, moves := parseInput(content)
	w.move(moves)
	return w.sumGps('O')
}

func (w warehouse) move2(moves []rune) {
	robotPos := w.robotPos()

	var shift func(pos point, move rune)
	shift = func(pos point, move rune) {
		step := movesMap[move]
		if w.at(pos) == '@' {
			shift(pos.add(step), move)
			swap(w, pos, pos.add(step))
		} else if w.at(pos) == '.' {
			return
		} else {
			if move == '<' || move == '>' {
				shift(pos.add(step).add(step), move)
				swap(w, pos.add(step), pos.add(step).add(step))
				swap(w, pos, pos.add(step))
			} else {
				var secondPos point
				if w.at(pos) == ']' {
					secondPos = pos.add(point{0, -1})
				} else {
					secondPos = pos.add(point{0, 1})
				}

				shift(pos.add(step), move)
				shift(secondPos.add(step), move)

				swap(w, pos, pos.add(step))
				swap(w, secondPos, secondPos.add(step))
			}
		}
	}

	var canBeMoved func([]point, rune) bool
	canBeMoved = func(positions []point, move rune) bool {
		step := movesMap[move]

		for _, pos := range positions {
			if w.at(pos) == '@' {
				return canBeMoved([]point{pos.add(step)}, move)
			} else if w.at(pos) == '[' || w.at(pos) == ']' {
				if move == '<' || move == '>' {
					return canBeMoved([]point{pos.add(step).add(step)}, move)
				} else {
					var secondPos point
					if w.at(pos) == ']' {
						secondPos = pos.add(point{0, -1})
					} else {
						secondPos = pos.add(point{0, 1})
					}

					newPositions := []point{pos.add(step), secondPos.add(step)}
					if !canBeMoved(newPositions, move) {
						return false
					}
				}

			} else if w.at(pos) == '.' {
				continue
			} else {
				return false
			}
		}

		return true
	}

	for _, move := range moves {
		step := movesMap[move]
		newPos := robotPos.add(step)

		if w.at(newPos) == '.' {
			swap(w, newPos, robotPos)
			robotPos = newPos
		} else if w.at(newPos) == '#' {
			continue
		} else {
			if canBeMoved([]point{robotPos}, move) {
				shift(robotPos, move)
				robotPos = newPos
			}
		}

	}
}

func part2(content string) int {
	w, moves := parseInput(content)
	_ = moves

	widenW := make(warehouse, len(w))
	transformations := map[rune]string{
		'#': "##",
		'O': "[]",
		'.': "..",
		'@': "@.",
	}

	for i, row := range w {
		widenW[i] = make([]rune, len(row)*2)
		for j, ch := range row {
			for k, newCh := range transformations[ch] {
				widenW[i][j*2+k] = newCh
			}
		}
	}

	widenW.move2(moves)
	return widenW.sumGps('[')
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
