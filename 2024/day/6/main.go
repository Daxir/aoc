package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	board, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	path, isLooping, err := findPath(board)
	if err != nil {
		fmt.Printf("error finding path: %v", err)
		return
	}
	if isLooping {
		fmt.Println("Path is looping")
		return
	}

	uniquePositions := make(map[coordinate]bool)
	for _, step := range path {
		uniquePositions[step.coordinate] = true
	}
	uniquePositionsCount := len(uniquePositions)
	fmt.Printf("(Part one) Unique positions visited: %v\n", uniquePositionsCount)

	obstructions := getPossibleObstructions(path)

	validObstructionCount := 0
	for _, obstruction := range obstructions {
		if evaluateObstruction(board, obstruction) {
			validObstructionCount++
		}
	}
	fmt.Printf("(Part two) Valid obstructions: %v\n", validObstructionCount)
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

type direction rune

const (
	up    direction = '^'
	down  direction = 'v'
	left  direction = '<'
	right direction = '>'
)

type coordinate struct {
	x int
	y int
}

func (c coordinate) move(d direction) coordinate {
	switch d {
	case up:
		return coordinate{x: c.x, y: c.y - 1}
	case down:
		return coordinate{x: c.x, y: c.y + 1}
	case left:
		return coordinate{x: c.x - 1, y: c.y}
	case right:
		return coordinate{x: c.x + 1, y: c.y}
	}
	return c
}

func rotateRight(dir direction) direction {
	switch dir {
	case up:
		return right
	case right:
		return down
	case down:
		return left
	case left:
		return up
	}
	return up
}

func (c coordinate) isOutOfBounds(board [][]rune) bool {
	return c.y < 0 || c.y >= len(board) || c.x < 0 || c.x >= len(board[0])
}

func findGuard(board [][]rune) (coordinate, error) {
	for y, row := range board {
		for x, cell := range row {
			if cell == rune(up) || cell == rune(down) || cell == rune(left) || cell == rune(right) {
				return coordinate{x: x, y: y}, nil
			}
		}
	}
	return coordinate{}, fmt.Errorf("guard not found")
}

type pathStep struct {
	coordinate
	direction
}

func findPath(board [][]rune) (path []pathStep, isLooping bool, err error) {
	boardCopy := make([][]rune, len(board))
	for i, row := range board {
		boardCopy[i] = make([]rune, len(row))
		copy(boardCopy[i], row)
	}
	guard, err := findGuard(boardCopy)
	if err != nil {
		return []pathStep{}, false, fmt.Errorf("error finding guard: %w", err)
	}

	visitedTurns := map[pathStep]bool{}

	path = []pathStep{}
	direction := direction(boardCopy[guard.y][guard.x])
	for {
		if guard.isOutOfBounds(boardCopy) {
			break
		}

		target := guard.move(direction)
		isTargetOutOfBounds := target.isOutOfBounds(boardCopy)

		if !isTargetOutOfBounds {
			targetCell := boardCopy[target.y][target.x]
			for targetCell == '#' {
				direction = rotateRight(direction)
				target = guard.move(direction)
				targetCell = boardCopy[target.y][target.x]

				if visitedTurns[pathStep{coordinate: target, direction: direction}] {
					return path, true, nil
				}
				visitedTurns[pathStep{coordinate: target, direction: direction}] = true
			}
		}

		path = append(path, pathStep{coordinate: guard, direction: direction})
		boardCopy[guard.y][guard.x] = '.'

		if isTargetOutOfBounds {
			break
		}

		boardCopy[target.y][target.x] = rune(direction)
		guard = target
	}

	return path, false, nil
}

func getPossibleObstructions(path []pathStep) []coordinate {
	uniquePositions := make(map[coordinate]bool)
	for _, step := range path {
		uniquePositions[step.coordinate] = true
	}

	obstructions := []coordinate{}
	for step := range uniquePositions {
		obstructions = append(obstructions, step)
	}

	return obstructions
}

func evaluateObstruction(board [][]rune, obstruction coordinate) bool {
	boardCopy := make([][]rune, len(board))
	for i, row := range board {
		boardCopy[i] = make([]rune, len(row))
		copy(boardCopy[i], row)
	}

	if boardCopy[obstruction.y][obstruction.x] != '.' {
		return false
	}

	boardCopy[obstruction.y][obstruction.x] = '#'

	_, isLooping, err := findPath(boardCopy)
	if err != nil {
		fmt.Printf("error finding path: %v", err)
		return false
	}

	return isLooping
}
