package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

type point [2]int
type keypad [][]rune

var nKeypad = keypad{
	{'7', '8', '9'},
	{'4', '5', '6'},
	{'1', '2', '3'},
	{' ', '0', 'A'},
}

var dKeypad = keypad{
	{' ', '^', 'A'},
	{'<', 'v', '>'},
}

var nStartPos = point{3, 2}
var dStartPos = point{0, 2}

func (k keypad) posMap() map[rune]point {
	result := map[rune]point{}
	for i, row := range k {
		for j := range row {
			result[k[i][j]] = point{i, j}
		}
	}
	return result
}

func moveHor(y1, y2 int) string {
	if dy := y2 - y1; dy > 0 {
		return strings.Repeat(">", dy)
	} else {
		return strings.Repeat("<", -dy)
	}
}

func moveVer(x1, x2 int) string {
	if dx := x2 - x1; dx > 0 {
		return strings.Repeat("v", dx)
	} else {
		return strings.Repeat("^", -dx)
	}
}

func move(kpPosMap map[rune]point, p1, p2 point) string {
	if p1 == p2 {
		return ""
	} else if p1[0] == p2[0] {
		return moveHor(p1[1], p2[1])
	} else if p1[1] == p2[1] {
		return moveVer(p1[0], p2[0])
	} else {
		gapPos := kpPosMap[' ']

		if p1[1] > p2[1] {
			// move left
			if p1[0] == gapPos[0] && p2[1] == gapPos[1] {
				// move vertically first if p1 is at the same level as gap
				return moveVer(p1[0], p2[0]) + moveHor(p1[1], p2[1])
			} else {
				// move horizontally first
				return moveHor(p1[1], p2[1]) + moveVer(p1[0], p2[0])
			}
		} else { // move right
			if p1[1] == gapPos[1] && p2[0] == gapPos[0] {
				// move horizontally first if p2 is at the same level as gap
				return moveHor(p1[1], p2[1]) + moveVer(p1[0], p2[0])
			} else {
				// move vertically first
				return moveVer(p1[0], p2[0]) + moveHor(p1[1], p2[1])
			}
		}
	}
}

func getPath(kpPosMap map[rune]point, code string, pos point) []string {
	result := make([]string, 0)
	currPos := pos

	for i := 0; i < len(code); i++ {
		nextPos := kpPosMap[rune(code[i])]
		result = append(result, move(kpPosMap, currPos, nextPos)+"A")
		currPos = nextPos
	}
	return result
}

func solve(codes []string, numRobots int) int {
	nPosMap := nKeypad.posMap()
	dPosMap := dKeypad.posMap()

	cache := map[string]int{}
	var getSequence func(code string, depth int) int
	getSequence = func(code string, depth int) int {
		if result, exists := cache[fmt.Sprintf("%s:%d", code, depth)]; exists {
			return result
		}

		var posMap map[rune]point
		var startPos point

		if unicode.IsNumber(rune(code[0])) {
			posMap = nPosMap
			startPos = nStartPos
		} else {
			posMap = dPosMap
			startPos = dStartPos
		}

		if depth == -1 {
			return len(code)

		} else {
			result := 0
			for _, move := range getPath(posMap, code, startPos) {
				result += getSequence(move, depth-1)
			}
			cache[fmt.Sprintf("%s:%d", code, depth)] = result
			return result
		}
	}

	result := 0
	for _, code := range codes {
		codeNum, _ := strconv.Atoi(code[:len(code)-1])
		lenSeq := getSequence(code, numRobots)
		result += codeNum * lenSeq
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
	codes := strings.Split(content, "\n")

	fmt.Println("part 1", solve(codes, 2))
	fmt.Println("part 2", solve(codes, 25))
}
