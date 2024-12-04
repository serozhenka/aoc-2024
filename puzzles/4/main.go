package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func part1(letters [][]rune) int {
	result := 0

	directions := [][2]int{
		{0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
	}
	lookup := []rune("XMAS")

	for i := 0; i < len(letters); i++ {
		for j := 0; j < len(letters[i]); j++ {
			for _, direction := range directions {
				stepI, stepJ := direction[0], direction[1]
				for lookupIdx, lookupLetter := range lookup {
					letterI := i + lookupIdx*stepI
					letterJ := j + lookupIdx*stepJ

					if letterI < 0 || letterI >= len(letters) ||
						letterJ < 0 || letterJ >= len(letters[i]) ||
						letters[letterI][letterJ] != lookupLetter {
						break
					}

					if lookupIdx == len(lookup)-1 {
						result += 1
					}
				}
			}
		}
	}

	return result
}

func part2(letters [][]rune) int {
	result := 0

	for i := 1; i < len(letters)-1; i++ {
		for j := 1; j < len(letters[i])-1; j++ {
			if letters[i][j] != 'A' {
				continue
			}

			oppositeLetters := map[rune]rune{'M': 'S', 'S': 'M'}
			if oppositeLetters[letters[i-1][j-1]] == letters[i+1][j+1] &&
				oppositeLetters[letters[i-1][j+1]] == letters[i+1][j-1] {
				result += 1
			}

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
	letters := make([][]rune, len(lines))
	for i, line := range strings.Split(content, "\n") {
		letters[i] = make([]rune, len(line))
		for j, letter := range line {
			letters[i][j] = letter
		}
	}

	fmt.Println("part 1", part1(letters))
	fmt.Println("part 2", part2(letters))
}
