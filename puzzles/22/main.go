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

type cappedQueue[T any] struct {
	capacity int
	elements []T
}

func makeCappedQueue[T any](capacity int) cappedQueue[T] {
	return cappedQueue[T]{capacity: capacity, elements: make([]T, 0, capacity)}
}

func (q *cappedQueue[T]) Size() int {
	return len(q.elements)
}

func (q *cappedQueue[T]) Dequeue() T {
	first := q.elements[0]
	q.elements = q.elements[1:]
	return first
}

func (q *cappedQueue[T]) Enqueue(el T) {
	if q.Size() >= q.capacity {
		q.Dequeue()
	}
	q.elements = append(q.elements, el)
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	contentBytes, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

	secrets := strings.Split(content, "\n")
	var secretHistory func(secret int, depth int) []int
	secretHistory = func(secret int, depth int) []int {
		if depth == 0 {
			return []int{}
		}

		history := []int{secret}
		step1 := ((secret * 64) ^ secret) % 16777216
		step2 := ((step1 / 32) ^ step1) % 16777216
		newSecret := ((step2 * 2048) ^ step2) % 16777216
		return append(history, secretHistory(newSecret, depth-1)...)
	}

	result := 0
	overallSequenceMap := map[[4]int]int{}

	for _, secretStr := range secrets {

		secret, _ := strconv.Atoi(secretStr)
		history := secretHistory(secret, 2001)
		result += history[len(history)-1]

		sequenceMap := map[[4]int]int{}

		sequence := makeCappedQueue[int](4)
		for i := 1; i < len(history); i++ {
			sequence.Enqueue((history[i] % 10) - (history[i-1] % 10))
			if sequence.Size() == 4 {
				var arr [4]int
				copy(arr[:], sequence.elements)

				if _, exists := sequenceMap[arr]; !exists {
					sequenceMap[arr] += history[i] % 10
				}
			}

		}

		for arr, bananas := range sequenceMap {
			overallSequenceMap[arr] += bananas
		}
	}

	fmt.Println("part 1", result)
	maxBananas := math.MinInt
	for _, bananas := range overallSequenceMap {
		maxBananas = max(maxBananas, bananas)
	}
	fmt.Println("part 2", maxBananas)
}
