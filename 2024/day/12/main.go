package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {
	garden, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	plots := mapToPlots(garden)
	regions := groupPlots(plots)

	totalPrice := 0
	for _, regions := range regions {
		for _, region := range regions {
			totalPrice += getRegionPrice(region)
		}
	}

	fmt.Printf("(Part one) Total price: %d\n", totalPrice)

	discountedPrice := 0
	for _, regions := range regions {
		for _, region := range regions {
			discountedPrice += getDiscountedRegionPrice(garden, region)
		}
	}

	fmt.Printf("(Part two) Discounted price: %d\n", discountedPrice)
}

func readInput() ([][]rune, error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	text := make([][]rune, 0)
	text = append(text, []rune{})
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		char := scanner.Text()
		if char == "\n" {
			text = append(text, []rune{})
			continue
		}
		charAsRunes := []rune(char)
		text[len(text)-1] = append(text[len(text)-1], charAsRunes[0])
	}
	if len(text[len(text)-1]) == 0 {
		text = text[:len(text)-1]
	}

	width := len(text[0])
	for _, line := range text {
		if len(line) != width {
			return nil, fmt.Errorf("all lines should have the same length")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return text, nil
}

type coordinate struct {
	x int
	y int
}

func (c coordinate) getAdjacent() []coordinate {
	up := coordinate{x: c.x - 1, y: c.y}
	down := coordinate{x: c.x + 1, y: c.y}
	left := coordinate{x: c.x, y: c.y - 1}
	right := coordinate{x: c.x, y: c.y + 1}

	return []coordinate{up, down, left, right}
}

func (c coordinate) isWithinBounds(garden [][]rune) bool {
	if c.x < 0 || c.x >= len(garden) {
		return false
	}

	if c.y < 0 || c.y >= len(garden[c.x]) {
		return false
	}

	return true
}

type plot struct {
	coordinate
	value     rune
	neighbors map[rune][]*plot
}

func mapToPlots(garden [][]rune) [][]*plot {
	nodes := make([][]*plot, 0)
	for y, line := range garden {
		nodes = append(nodes, make([]*plot, 0))
		for x, cell := range line {
			nodes[y] = append(nodes[y], &plot{
				coordinate: coordinate{x: x, y: y},
				value:      cell,
				neighbors:  make(map[rune][]*plot),
			})
		}
	}

	for _, line := range nodes {
		for _, n := range line {
			adjacents := n.getAdjacent()
			for _, adj := range adjacents {
				if !adj.isWithinBounds(garden) {
					continue
				}

				adjNode := nodes[adj.y][adj.x]
				n.neighbors[adjNode.value] = append(n.neighbors[adjNode.value], adjNode)
			}
		}
	}

	return nodes
}

func dfs(current *plot, visited map[*plot]bool, group *[]*plot, targetValue rune) {
	visited[current] = true
	*group = append(*group, current)

	for _, neighbor := range current.neighbors[targetValue] {
		if !visited[neighbor] {
			dfs(neighbor, visited, group, targetValue)
		}
	}
}

func groupPlots(plots [][]*plot) map[rune][][]*plot {
	visited := make(map[*plot]bool)
	groups := make(map[rune][][]*plot)

	for _, row := range plots {
		for _, p := range row {
			if !visited[p] {
				var group []*plot
				dfs(p, visited, &group, p.value)
				groups[p.value] = append(groups[p.value], group)
			}
		}
	}

	return groups
}

func getRegionPrice(region []*plot) int {
	area := len(region)

	perimeter := 0
	for _, plot := range region {
		perimeter += 4 - len(plot.neighbors[plot.value])
	}

	return area * perimeter
}

// This is the worst absolute solution I could come up with
// I hate it, but I hate the problem even more
// This stays here unless I magically stop hating the problem and come up with a better solution.
func getDiscountedRegionPrice(garden [][]rune, region []*plot) int {
	area := len(region)

	sides := scanLeft(garden, region)
	sides += scanRight(garden, region)
	sides += scanTop(garden, region)
	sides += scanBottom(garden, region)

	return area * sides
}

func countLines(line []int) int {
	sorted := slices.Sorted(slices.Values(line))

	count := 1
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1]+1 {
			count++
		}
	}

	return count
}

func (c coordinate) getLeft() coordinate {
	return coordinate{x: c.x - 1, y: c.y}
}

func (c coordinate) getRight() coordinate {
	return coordinate{x: c.x + 1, y: c.y}
}

func (c coordinate) getUp() coordinate {
	return coordinate{x: c.x, y: c.y - 1}
}

func (c coordinate) getDown() coordinate {
	return coordinate{x: c.x, y: c.y + 1}
}

func scanLeft(garden [][]rune, region []*plot) int {
	potentialLines := make(map[int][]int)
	for _, plot := range region {
		if len(plot.neighbors[plot.value]) >= 4 {
			continue
		}

		left := plot.getLeft()
		if !left.isWithinBounds(garden) {
			potentialLines[left.x] = append(potentialLines[left.x], left.y)
			continue
		}

		leftPlot := garden[left.y][left.x]
		if leftPlot != plot.value {
			potentialLines[left.x] = append(potentialLines[left.x], left.y)
		}
	}

	lines := 0
	for _, line := range potentialLines {
		lines += countLines(line)
	}

	return lines
}

func scanRight(garden [][]rune, region []*plot) int {
	potentialLines := make(map[int][]int)
	for _, plot := range region {
		if len(plot.neighbors[plot.value]) >= 4 {
			continue
		}

		right := plot.getRight()
		if !right.isWithinBounds(garden) {
			potentialLines[right.x] = append(potentialLines[right.x], right.y)
			continue
		}

		rightPlot := garden[right.y][right.x]
		if rightPlot != plot.value {
			potentialLines[right.x] = append(potentialLines[right.x], right.y)
		}
	}

	lines := 0
	for _, line := range potentialLines {
		lines += countLines(line)
	}

	return lines
}

func scanTop(garden [][]rune, region []*plot) int {
	potentialLines := make(map[int][]int)
	for _, plot := range region {
		if len(plot.neighbors[plot.value]) >= 4 {
			continue
		}

		top := plot.getUp()
		if !top.isWithinBounds(garden) {
			potentialLines[top.y] = append(potentialLines[top.y], top.x)
			continue
		}

		topPlot := garden[top.y][top.x]
		if topPlot != plot.value {
			potentialLines[top.y] = append(potentialLines[top.y], top.x)
		}
	}

	lines := 0
	for _, line := range potentialLines {
		lines += countLines(line)
	}

	return lines
}

func scanBottom(garden [][]rune, region []*plot) int {
	potentialLines := make(map[int][]int)
	for _, plot := range region {
		if len(plot.neighbors[plot.value]) >= 4 {
			continue
		}

		bottom := plot.getDown()
		if !bottom.isWithinBounds(garden) {
			potentialLines[bottom.y] = append(potentialLines[bottom.y], bottom.x)
			continue
		}

		bottomPlot := garden[bottom.y][bottom.x]
		if bottomPlot != plot.value {
			potentialLines[bottom.y] = append(potentialLines[bottom.y], bottom.x)
		}
	}

	lines := 0
	for _, line := range potentialLines {
		lines += countLines(line)
	}

	return lines
}
