package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	layout, err := readInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	antennaInfo := getAntennaInfo(layout)

	antinodes := getResonantAntinodes(antennaInfo, layout)
	uniqueResonantAntinodeCount := getUniqueAntinodeCount(antinodes)
	fmt.Printf("(Part one) Unique antinodes: %v\n", uniqueResonantAntinodeCount)

	linearAntiNodes := getLinearAntinodes(antennaInfo, layout)
	uniqueLinearAntinodeCount := getUniqueAntinodeCount(linearAntiNodes)
	fmt.Printf("(Part two) Unique linear antinodes: %v\n", uniqueLinearAntinodeCount)
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

func (c coordinate) findOtherEnd(center coordinate) coordinate {
	x := center.x + (center.x - c.x)
	y := center.y + (center.y - c.y)
	return coordinate{x: x, y: y}
}

func (c coordinate) isWithinBounds(layout [][]rune) bool {
	if c.x < 0 || c.y < 0 {
		return false
	}
	if c.y >= len(layout) || c.x >= len(layout[0]) {
		return false
	}
	return true
}

func getAntennaInfo(layout [][]rune) map[rune][]coordinate {
	antennas := make(map[rune][]coordinate)
	for y, line := range layout {
		for x, char := range line {
			if char == '.' {
				continue
			}
			antennas[char] = append(antennas[char], coordinate{x: x, y: y})
		}
	}
	return antennas
}

func getAllPairs(coordinates []coordinate) [][]coordinate {
	pairs := make([][]coordinate, 0)
	for i, c1 := range coordinates {
		for j, c2 := range coordinates {
			if i == j {
				continue
			}
			pairs = append(pairs, []coordinate{c1, c2})
		}
	}
	return pairs
}

func getResonantAntinodes(antennaInfo map[rune][]coordinate, layout [][]rune) []coordinate {
	antinodes := []coordinate{}
	for _, coordinates := range antennaInfo {
		pairs := getAllPairs(coordinates)
		for _, pair := range pairs {
			antinode := pair[0].findOtherEnd(pair[1])
			if antinode.isWithinBounds(layout) {
				antinodes = append(antinodes, antinode)
			}
		}
	}
	return antinodes
}

func getLinearAntinodes(antennaInfo map[rune][]coordinate, layout [][]rune) []coordinate {
	linearAntiNodes := []coordinate{}
	for _, coordinates := range antennaInfo {
		pairs := getAllPairs(coordinates)
		for _, pair := range pairs {
			linearAntiNodes = append(linearAntiNodes, pair[0])

			antinode := pair[0].findOtherEnd(pair[1])
			if antinode.isWithinBounds(layout) {
				layout[antinode.y][antinode.x] = '#'
				linearAntiNodes = append(linearAntiNodes, antinode)

				from := pair[1]
				center := antinode
				for {
					next := from.findOtherEnd(center)
					if !next.isWithinBounds(layout) {
						break
					}
					linearAntiNodes = append(linearAntiNodes, next)
					from = center
					center = next
				}
			}
		}
	}

	return linearAntiNodes
}

func getUniqueAntinodeCount(antinodes []coordinate) int {
	uniqueAntinodes := make(map[coordinate]bool)
	for _, antinode := range antinodes {
		uniqueAntinodes[antinode] = true
	}
	return len(uniqueAntinodes)
}
