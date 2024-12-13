package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

var steps = [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

type point [2]int

type plot struct {
	value   rune
	tlpoint point // top-left point
}

type region []plot

func isStraightLine(p1, p2, p3 point) bool {
	return (p1[0] == p2[0] && p2[0] == p3[0]) || (p1[1] == p2[1] && p2[1] == p3[1])
}

func distance(p1, p2 point) float64 {
	return math.Sqrt(math.Pow(float64(p1[0]-p2[0]), 2) + math.Pow(float64(p1[1]-p2[1]), 2))
}

func (p1 point) add(p2 point) point {
	return point{p1[0] + p2[0], p1[1] + p2[1]}
}

func (p plot) equal(other plot) bool {
	return p.tlpoint == other.tlpoint
}

func (p plot) in(lst region) bool {
	for _, other := range lst {
		if p.tlpoint == other.tlpoint {
			return true
		}
	}
	return false
}

func (p plot) isNeighbor(other plot) bool {
	if p.value != other.value {
		return false
	}

	for _, step := range steps {
		if p.tlpoint == other.tlpoint.add(step) {
			return true
		}
	}

	return false
}

func getPlot(plotsMap [][]plot, p point) (plot, bool) {
	for _, row := range plotsMap {
		for _, plt := range row {
			if plt.tlpoint == p {
				return plt, true
			}
		}
	}
	return plot{}, false
}

func (p1 point) isConnectedTo(p2 point, r region, plotsMap [][]plot) bool {
	value := r[0].value
	if distance(p1, p2) != 1 {
		return false
	}

	var plot1, plot2 plot
	if p1[0]-p2[0] == 0 { // horizontal
		var left point
		if p1[1] < p2[1] {
			left = p1
		} else {
			left = p2
		}

		plot1 = plot{value: value, tlpoint: left}
		plot2 = plot{value: value, tlpoint: left.add(point{-1, 0})}
	} else { // vertical
		var top point
		if p1[0] < p2[0] {
			top = p1
		} else {
			top = p2
		}

		plot1 = plot{value: value, tlpoint: top}
		plot2 = plot{value: value, tlpoint: top.add(point{0, -1})}
	}

	rp1, suc1 := getPlot(plotsMap, plot1.tlpoint)
	rp2, suc2 := getPlot(plotsMap, plot2.tlpoint)
	return (!(plot1.in(r) && plot2.in(r))) &&
		((suc1 && rp1.value == value) || (suc2 && rp2.value == value))
}

func (r region) sides(plotsMap [][]plot) int {
	points := []point{}
	var steps = []point{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	for _, plt := range r {
		for _, step := range steps {
			newPoint := plt.tlpoint.add(step)
			if !slices.Contains(points, newPoint) {
				points = append(points, plt.tlpoint.add(step))
			}
		}
	}

	connections := map[point][]point{}
	for _, p1 := range points {
		pointConnections := make([]point, 0, 2)
		for _, p2 := range points {
			if p1 == p2 {
				continue
			}

			if p1.isConnectedTo(p2, r, plotsMap) {
				pointConnections = append(pointConnections, p2)
			}
		}

		connections[p1] = pointConnections
	}

	numCorners := 0
	for p, pConns := range connections {
		if len(pConns) == 4 {
			numCorners += 1
		} else if len(pConns) != 2 {
			continue
		}
		p2, p3 := pConns[0], pConns[1]
		if !isStraightLine(p, p2, p3) {
			numCorners += 1
		}
	}

	return numCorners
}

func (r region) area() int {
	return len(r)
}

func (r region) perimeter() int {
	allNeighborsCount := 0
	for _, plt := range r {
		for _, otherPlot := range r {
			if plt.equal(otherPlot) {
				continue
			}

			if plt.isNeighbor(otherPlot) {
				allNeighborsCount += 1
			}
		}
	}

	return len(r)*4 - allNeighborsCount
}

func (r region) price() int {
	return r.area() * r.perimeter()
}

func (r region) price2(plotsMap [][]plot) int {
	return r.area() * r.sides(plotsMap)
}

func getNeighborPlots(plt plot, plotsMap [][]plot, r region) region {
	r = append(r, plt)
	for _, row := range plotsMap {
		for _, mbNeighbor := range row {
			if plt.equal(mbNeighbor) || !plt.isNeighbor(mbNeighbor) || mbNeighbor.in(r) {
				continue
			}

			r = append(r, mbNeighbor)
			for _, newNeighbor := range getNeighborPlots(mbNeighbor, plotsMap, r) {
				if !newNeighbor.in(r) {
					r = append(r, newNeighbor)
				}
			}
		}
	}
	return r
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

	plotsMap := make([][]plot, len(lines))
	for i, line := range lines {
		plotsMap[i] = make([]plot, len(line))
		for j, ch := range line {
			plotsMap[i][j] = plot{value: ch, tlpoint: [2]int{i, j}}
		}
	}

	regions := make([]region, 0)
	for _, row := range plotsMap {
		for _, plt := range row {
			alreadyInRegion := false
			for _, r := range regions {
				if plt.in(r) {
					alreadyInRegion = true
					break
				}
			}
			if alreadyInRegion {
				continue
			}

			regions = append(regions, getNeighborPlots(plt, plotsMap, region{}))
		}
	}

	totalPrice := 0
	totalPrice2 := 0
	for _, r := range regions {
		totalPrice += r.price()
		totalPrice2 += r.price2(plotsMap)
	}

	fmt.Println("part 1", totalPrice)
	fmt.Println("part 2", totalPrice2)
}
