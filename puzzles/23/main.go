package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

func filter[T any](s []T, f func(T) bool) []T {
	result := make([]T, 0)
	for _, el := range s {
		if f(el) {
			result = append(result, el)
		}
	}
	return result
}

func part1(adjacencyMap map[string]map[string]bool) int {
	networks := make([][3]string, 0)
	visited := map[string]bool{}

	for v1 := range adjacencyMap {
		innerVisited := map[string]bool{}
		for v2 := range adjacencyMap[v1] {
			if visited[v2] {
				continue
			}

			for v3 := range adjacencyMap[v2] {
				if v1 == v3 || visited[v3] || innerVisited[v3] {
					continue
				}

				if adjacencyMap[v1][v3] {
					networks = append(networks, [3]string{v1, v2, v3})
				}
			}

			innerVisited[v2] = true
		}
		visited[v1] = true
	}

	result := 0
	for _, n := range networks {
		for _, v := range n {
			if strings.HasPrefix(v, "t") {
				result += 1
				break
			}
		}
	}

	return result
}

func part2(adjacencyMap map[string]map[string]bool) string {
	var bronKerbosh func(r, p, x []string) [][]string
	bronKerbosh = func(r, p, x []string) [][]string {
		if len(p) == 0 && len(x) == 0 {
			return [][]string{r}
		}

		maxCliques := make([][]string, 0)

		for _, v := range p {
			for _, n := range bronKerbosh(
				append(r, v),
				filter(p, func(el string) bool {
					return adjacencyMap[v][el]
				}),
				filter(x, func(el string) bool {
					return adjacencyMap[v][el]
				})) {
				nCopy := make([]string, len(n))
				copy(nCopy, n)
				maxCliques = append(maxCliques, nCopy)
			}

			p = slices.DeleteFunc(p, func(el string) bool {
				return el == v
			})
			x = append(x, v)
		}

		return maxCliques
	}

	vertices := make([]string, 0, len(adjacencyMap))
	for v := range adjacencyMap {
		vertices = append(vertices, v)
	}

	maxCliques := bronKerbosh([]string{}, vertices, []string{})
	longestClique := []string{}
	for _, c := range maxCliques {
		if len(c) > len(longestClique) {
			longestClique = c
		}
	}
	slices.Sort(longestClique)
	return strings.Join(longestClique, ",")
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

	edges := make([][2]string, 0)
	uniqueVertices := map[string]bool{}

	for _, line := range lines {
		v1, v2, _ := strings.Cut(line, "-")
		edges = append(edges, [2]string{v1, v2})
		uniqueVertices[v1] = true
		uniqueVertices[v2] = true
	}

	adjacencyMap := map[string]map[string]bool{}
	for v := range uniqueVertices {
		adjacencyMap[v] = map[string]bool{}
	}

	for _, edge := range edges {
		v1, v2 := edge[0], edge[1]
		adjacencyMap[v1][v2] = true
		adjacencyMap[v2][v1] = true
	}

	fmt.Println("part 1", part1(adjacencyMap))
	fmt.Println("part 2", part2(adjacencyMap))
}
