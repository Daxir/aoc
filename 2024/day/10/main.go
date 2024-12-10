package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	tMap, err := readInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	nodes := mapToNodes(tMap)
	zeroNodes := findValueNodes(nodes, 0)

	uniqueSum := 0
	nonUniqueSum := 0
	for _, zeroNode := range zeroNodes {
		paths := collectPaths(zeroNode, 0)
		uniquePaths := make(map[coordinate]bool)
		for _, path := range paths {
			uniquePaths[path.coordinate] = true
		}
		uniqueSum += len(uniquePaths)
		nonUniqueSum += len(paths)
	}

	fmt.Printf("(Part one) Unique sum: %d\n", uniqueSum)
	fmt.Printf("(Part two) Non-unique sum: %d\n", nonUniqueSum)
}

var digits = map[string]int{
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
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

func (c coordinate) isWithinBounds(tMap [][]int) bool {
	if c.x < 0 || c.x >= len(tMap) {
		return false
	}

	if c.y < 0 || c.y >= len(tMap[c.x]) {
		return false
	}

	return true
}

func readInput() ([][]int, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	tMap := make([][]int, 0)
	tMap = append(tMap, make([]int, 0))
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		char := scanner.Text()
		if char == "\n" {
			tMap = append(tMap, make([]int, 0))
			continue
		}
		digit, ok := digits[char]
		if !ok {
			return nil, fmt.Errorf("invalid digit: %s", char)
		}
		tMap[len(tMap)-1] = append(tMap[len(tMap)-1], digit)
	}
	if len(tMap[len(tMap)-1]) == 0 {
		tMap = tMap[:len(tMap)-1]
	}

	width := len(tMap[0])
	for _, row := range tMap {
		if len(row) != width {
			return nil, fmt.Errorf("all rows should have the same length: %d", len(row))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return tMap, nil
}

type node struct {
	coordinate
	value     int
	neighbors map[int][]*node
}

func mapToNodes(tMap [][]int) [][]*node {
	nodes := make([][]*node, 0)
	for y, line := range tMap {
		nodes = append(nodes, make([]*node, 0))
		for x, cell := range line {
			nodes[y] = append(nodes[y], &node{
				coordinate: coordinate{x: x, y: y},
				value:      cell,
				neighbors:  make(map[int][]*node),
			})
		}
	}

	for _, line := range nodes {
		for _, n := range line {
			adjacents := n.getAdjacent()
			for _, adj := range adjacents {
				if !adj.isWithinBounds(tMap) {
					continue
				}

				adjNode := nodes[adj.y][adj.x]
				n.neighbors[adjNode.value] = append(n.neighbors[adjNode.value], adjNode)
			}
		}
	}

	return nodes
}

func findValueNodes(nodes [][]*node, value int) []*node {
	valueNodes := make([]*node, 0)
	for _, line := range nodes {
		for _, n := range line {
			if n.value == value {
				valueNodes = append(valueNodes, n)
			}
		}
	}

	return valueNodes
}

func collectPaths(current *node, targetValue int) []*node {
	if current.value == targetValue && targetValue == 9 {
		return []*node{current}
	}

	if current.value != targetValue {
		return nil
	}

	nextValue := targetValue + 1
	var result []*node
	for _, neighbor := range current.neighbors[nextValue] {
		result = append(result, collectPaths(neighbor, nextValue)...)
	}

	return result
}
