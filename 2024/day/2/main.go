package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	reports, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	solution := solution{reports: reports}

	validCount := solution.solve(isValid)

	fmt.Printf("(Part one) valid count: %v\n", validCount)

	validCount = solution.solve(isValidWithTolerance)
	fmt.Printf("(Part two) valid count: %v\n", validCount)
}

func readInput() (reports [][]int, err error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`\s+`)
		split := re.Split(line, -1)

		report := make([]int, 0)
		for _, s := range split {
			value, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("error converting to int: %w", err)
			}
			report = append(report, value)
		}
		reports = append(reports, report)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return reports, nil
}

type solution struct {
	reports [][]int
}

func (p solution) solve(isValidFunc func([]int) bool) int {
	counter := 0
	for _, report := range p.reports {
		if isValidFunc(report) {
			counter++
		}
	}
	return counter
}

func isValid(report []int) bool {
	safeIncreases := []int{-3, -2, -1, 1, 2, 3}

	if !isDecreasingOrIncreasing(report) {
		return false
	}

	for i := 0; i < len(report)-1; i++ {
		increase := report[i] - report[i+1]
		if !slices.Contains(safeIncreases, increase) {
			return false
		}
	}

	return true
}

func isDecreasingOrIncreasing(report []int) bool {
	isAscending := slices.IsSorted(report)
	isDescending := slices.IsSortedFunc(report, func(i, j int) int { return j - i })
	return isAscending || isDescending
}

// a terrible brute force solution, I'm not proud of it.
func isValidWithTolerance(report []int) bool {
	if isValid(report) {
		return true
	}

	for i := 0; i < len(report); i++ {
		subslice := make([]int, 0)
		for j := 0; j < len(report); j++ {
			if j == i {
				continue
			}
			subslice = append(subslice, report[j])
		}
		if isValid(subslice) {
			return true
		}
	}

	return false
}
