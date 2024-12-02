package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	left, right, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	sum, err := partOne(left, right)
	if err != nil {
		fmt.Printf("error calculating part one: %v", err)
		return
	}

	fmt.Printf("(Part one) sum of distances: %v\n", int(sum))

	similarity := partTwo(left, right)
	fmt.Printf("(Part two) similarity: %v\n", similarity)
}

func readInput() (left, right []int, err error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	leftSlice := make([]int, 0)
	rightSlice := make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`\s+`)
		split := re.Split(line, -1)

		leftValue, err := strconv.Atoi(split[0])
		if err != nil {
			return nil, nil, fmt.Errorf("error converting to int: %v", err)
		}
		rightValue, err := strconv.Atoi(split[1])
		if err != nil {
			return nil, nil, fmt.Errorf("error converting to int: %v", err)
		}
		leftSlice = append(leftSlice, leftValue)
		rightSlice = append(rightSlice, rightValue)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error scanning file: %v", err)
	}

	return leftSlice, rightSlice, nil
}

func partOne(left, right []int) (int, error) {
	if len(left) != len(right) {
		fmt.Printf("left and right slices are not the same length")
		return 0, fmt.Errorf("left and right slices are not the same length")
	}

	slices.Sort(left)
	slices.Sort(right)

	distances := make([]float64, len(left))
	for i := 0; i < len(left); i++ {
		distances[i] = math.Abs(float64(left[i]) - float64(right[i]))
	}

	sum := 0.0
	for _, distance := range distances {
		sum += distance
	}

	return int(sum), nil
}

func partTwo(left, right []int) int {
	lookup := make(map[int]int)
	for _, value := range right {
		lookup[value]++
	}

	similarity := 0
	for _, value := range left {
		if lookup[value] > 0 {
			similarity += value * lookup[value]
		}
	}

	return similarity
}