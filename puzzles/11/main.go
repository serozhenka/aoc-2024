package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func blinkStone(stone int) []int {
	newStones := make([]int, 0, 2)
	if stone == 0 {
		newStones = append(newStones, 1)
	} else if stoneStr := strconv.Itoa(stone); len(stoneStr)%2 == 0 {
		firstStone, _ := strconv.Atoi(stoneStr[:len(stoneStr)/2])
		secondStone, _ := strconv.Atoi(stoneStr[len(stoneStr)/2:])
		newStones = append(newStones, firstStone, secondStone)
	} else {
		newStones = append(newStones, stone*2024)
	}

	return newStones
}

func blink(stones []int, times int) int {
	cache := map[[2]int]int{}

	var do func(stones []int, times int) int
	do = func(stones []int, times int) int {
		if times == 0 {
			return len(stones)
		}

		result := 0
		for _, stone := range stones {
			var newStones int
			cacheKey := [2]int{stone, times}
			if _, exists := cache[cacheKey]; exists {
				newStones = cache[cacheKey]
			} else {
				newStones = do(blinkStone(stone), times-1)
				cache[cacheKey] = newStones
			}
			result += newStones
		}
		return result
	}

	return do(stones, times)
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)
	stonesStrList := strings.Split(content, " ")

	stones := make([]int, len(stonesStrList))
	for i, stoneStr := range stonesStrList {
		stones[i], _ = strconv.Atoi(stoneStr)
	}

	fmt.Println("part 1", blink(stones, 25))
	fmt.Println("part 2", blink(stones, 75))
}
