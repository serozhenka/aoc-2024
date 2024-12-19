package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func parseInput(content string) ([]string, []string) {
	towelsStr, designsStr, _ := strings.Cut(content, "\n\n")
	return strings.Split(towelsStr, ", "), strings.Split(designsStr, "\n")
}

func countWays(design string, towels []string) int {
	cache := map[string]int{}
	var do func(design string) int
	do = func(design string) int {
		if _, exists := cache[design]; exists {
			return cache[design]
		}

		result := 0

		if design == "" {
			result += 1
		} else {
			for _, towel := range towels {
				if prefix, found := strings.CutPrefix(design, towel); found {
					result += do(prefix)
				}
			}
		}

		cache[design] = result
		return result
	}

	return do(design)
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)
	towels, designs := parseInput(content)

	var possibleDesigns, totalWays int
	for _, design := range designs {
		ways := countWays(design, towels)
		totalWays += ways
		if ways != 0 {
			possibleDesigns += 1
		}
	}

	fmt.Println("part 1", possibleDesigns)
	fmt.Println("part 2", totalWays)
}
