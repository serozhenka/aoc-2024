package main

import (
	"fmt"
	"log"
	"strings"

	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

func getBlocks(diskMap string) []string {
	blocks := make([]string, 0)
	for i, ch := range diskMap {
		count, _ := strconv.Atoi(string(ch))
		var newPart string
		if i%2 == 0 {
			newPart = strconv.Itoa(i / 2)
		} else {
			newPart = "."
		}

		for j := 0; j < count; j++ {
			blocks = append(blocks, newPart)
		}
	}

	return blocks
}

func checksum(blocks []string) int {
	result := 0
	for i, block := range blocks {
		if block != "." {
			num, _ := strconv.Atoi(block)
			result += num * i
		}
	}
	return result
}

func swap[T any](lst []T, idx1 int, idx2 int) {
	lst[idx1], lst[idx2] = lst[idx2], lst[idx1]
}

func part1(diskMap string) int {
	blocks := getBlocks(diskMap)

	freeIdx := 0
	for i := len(blocks) - 1; i >= freeIdx; {
		if blocks[freeIdx] != "." {
			freeIdx += 1
			continue
		}

		if blocks[i] == "." {
			i--
			continue
		}

		swap(blocks, i, freeIdx)
	}

	return checksum(blocks)
}

func count[T string](lst []T, target T) int {
	result := 0
	for _, item := range lst {
		if item == target {
			result += 1
		}
	}
	return result
}

func part2(diskMap string) int {
	blocks := getBlocks(diskMap)
	seen := map[string]int{}

	for i := len(blocks) - 1; i >= 0; {
		if blocks[i] == strings.Repeat(".", len(blocks[i])) || seen[blocks[i]] == 1 {
			i--
			continue
		}

		blockCount := count(blocks, blocks[i])
		seen[blocks[i]] = 1

		freeLength := 0
		for j := 0; j <= i; j++ {
			if blocks[j] == "." {
				freeLength += 1
			} else {
				freeLength = 0
			}

			if freeLength == blockCount {
				for k := 1; k <= blockCount; k++ {
					swap(blocks, j-blockCount+k, i-blockCount+k)
				}
				break
			}
		}

		i -= blockCount
	}

	return checksum(blocks)
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	diskMap := string(contentBytes)

	fmt.Println("part 1", part1(diskMap))
	fmt.Println("part 2", part2(diskMap))
}
