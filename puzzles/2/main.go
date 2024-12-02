package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

func strLstToIntLst(lst []string) []int {
	result := make([]int, len(lst))
	for i, s := range lst {
		val, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Error converting string to int: %v\n", err)
		}
		result[i] = val
	}
	return result
}

func isReportSafe(report []int, numBadLevelsAllowed int) bool {
	do := func(report []int, ascending bool) bool {
		reportCopy := append([]int(nil), report...)
		badLevels := 0

		for i := 1; i < len(reportCopy); {
			diff := reportCopy[i] - reportCopy[i-1]
			if (!ascending && !(diff >= 1 && diff <= 3)) ||
				(ascending && !(diff <= -1 && diff >= -3)) {
				badLevels += 1
				reportCopy = append(reportCopy[:i], reportCopy[i+1:]...)

				if badLevels > numBadLevelsAllowed {
					return false
				}
			} else {
				i += 1
			}
		}
		return true
	}

	reversedReport := append([]int(nil), report...)
	slices.Reverse(reversedReport)

	for _, order := range []bool{false, true} {
		if do(report, order) || do(reversedReport, order) {
			return true
		}
	}
	return false

}

func getSafeReportsCount(reports [][]int, numBadLevelsAllowed int) int {
	safeReportsCount := 0
	for _, report := range reports {
		if isReportSafe(report, numBadLevelsAllowed) {
			safeReportsCount += 1
		}
	}
	return safeReportsCount
}

func main() {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	content, err := os.ReadFile(filepath.Join(currentDir, "input.txt"))
	if err != nil {
		log.Fatal(err)
	}

	reports := make([][]int, 0)

	for _, line := range strings.Split(string(content), "\n") {
		splitted := strings.Split(line, " ")
		reports = append(reports, strLstToIntLst(splitted))
	}

	fmt.Println("part 1", getSafeReportsCount(reports, 0))
	fmt.Println("part 2", getSafeReportsCount(reports, 1))
}
