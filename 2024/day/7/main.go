package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	equations, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	result := getCalibrationResult(equations, []operator{add, mul})
	fmt.Printf("(Part one) Total calibration result: %d\n", result)

	result = getCalibrationResult(equations, []operator{add, mul, concat})
	fmt.Printf("(Part two) Total calibration result: %d\n", result)
}

type equation struct {
	result  int
	numbers []int
}

func readInput() ([]equation, error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	equations := make([]equation, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		eq := equation{}

		equation := strings.Split(line, ": ")
		resultValue, err := strconv.Atoi(equation[0])
		if err != nil {
			return nil, fmt.Errorf("error converting result value to int: %w", err)
		}

		eq.result = resultValue

		numbers := strings.Split(equation[1], " ")
		eq.numbers = make([]int, len(numbers))

		for i, number := range numbers {
			num, err := strconv.Atoi(number)
			if err != nil {
				return nil, fmt.Errorf("error converting number to int: %w", err)
			}
			eq.numbers[i] = num
		}
		equations = append(equations, eq)
	}

	return equations, nil
}

type operator string

const (
	add    operator = "+"
	mul    operator = "*"
	concat operator = "||"
)

func apply(a int, op operator, b int) int {
	switch op {
	case add:
		return a + b
	case mul:
		return a * b
	case concat:
		multiplier := 1
		for temp := b; temp > 0; temp /= 10 {
			multiplier *= 10
		}
		return a*multiplier + b
	}
	panic("unknown operator")
}

func isEquationPossible(eq equation, operators []operator) bool {
	if len(eq.numbers) == 1 {
		return eq.numbers[0] == eq.result
	}

	for _, operator := range operators {
		operationResult := apply(eq.numbers[0], operator, eq.numbers[1])
		newNumbers := make([]int, len(eq.numbers)-1)
		newNumbers[0] = operationResult
		copy(newNumbers[1:], eq.numbers[2:])
		if isEquationPossible(equation{
			result:  eq.result,
			numbers: newNumbers,
		}, operators) {
			return true
		}
	}

	return false
}

func getCalibrationResult(equations []equation, operators []operator) int {
	result := 0
	for _, eq := range equations {
		if isEquationPossible(eq, operators) {
			result += eq.result
		}
	}

	return result
}
