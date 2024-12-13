package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var buttonRegex = regexp.MustCompile(`Button .: X\+(\d+), Y\+(\d+)`)
var prizeRegex = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

func parseLine(regex *regexp.Regexp, input string) [2]float64 {
	match := regex.FindStringSubmatch(input)
	x, _ := strconv.Atoi(match[1])
	y, _ := strconv.Atoi(match[2])
	return [2]float64{float64(x), float64(y)}
}

func isFloatInteger(x float64) bool {
	rounded := math.Round(x)
	return math.Abs(x-rounded) <= math.Pow(10, -3)
}

func solveEquation(a [2]float64, b [2]float64, prize [2]float64) (float64, float64) {
	stepsB := (a[1]/a[0]*prize[0] - prize[1]) / (a[1]*b[0]/a[0] - b[1])
	stepsA := (prize[0] - b[0]*stepsB) / a[0]
	return stepsA, stepsB
}

func solve(input string, prizeIncrement float64) float64 {
	machinesStr := strings.Split(input, "\n\n")
	result := 0.

	for _, machineStr := range machinesStr {
		lines := strings.Split(machineStr, "\n")
		buttonA := parseLine(buttonRegex, lines[0])
		buttonB := parseLine(buttonRegex, lines[1])
		prize := parseLine(prizeRegex, lines[2])
		prize = [2]float64{prize[0] + prizeIncrement, prize[1] + prizeIncrement}

		stepsA, stepsB := solveEquation(buttonA, buttonB, prize)
		if !isFloatInteger(stepsA) || !isFloatInteger(stepsB) {
			continue
		}

		result += stepsA*3 + stepsB
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

	fmt.Println("part 1", fmt.Sprintf("%.f", solve(content, 0)))
	fmt.Println("part 2", fmt.Sprintf("%.f", solve(content, 10000000000000)))
}
