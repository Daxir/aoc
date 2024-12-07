package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mowshon/iterium"
)

func main() {
	equations, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	result, err := getCalibrationResult(equations, []operator{add, mul})
	if err != nil {
		fmt.Printf("error getting calibration result: %v", err)
		return
	}

	fmt.Printf("(Part one) Total calibration result: %d\n", result)

	result, err = getCalibrationResult(equations, []operator{add, mul, concat})
	if err != nil {
		fmt.Printf("error getting calibration result: %v", err)
		return
	}

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

func getPossibleOperators(eq equation, operators []operator) ([][]operator, error) {
	product := iterium.Product(operators, len(eq.numbers)-1)
	possibleOperators, err := product.Slice()
	if err != nil {
		return nil, fmt.Errorf("error getting possible operators: %w", err)
	}

	return possibleOperators, nil
}

func isEquationPossible(eq equation, operators []operator) (bool, error) {
	possibleOperators, err := getPossibleOperators(eq, operators)
	if err != nil {
		return false, fmt.Errorf("error getting possible operators: %w", err)
	}

	for _, operators := range possibleOperators {
		result := eq.numbers[0]
		for i, operator := range operators {
			if i+1 >= len(eq.numbers) {
				break
			}

			result = apply(result, operator, eq.numbers[i+1])
		}
		if result == eq.result {
			return true, nil
		}
	}

	return false, nil
}

func getCalibrationResult(equations []equation, operators []operator) (int, error) {
	result := 0
	for _, eq := range equations {
		possible, err := isEquationPossible(eq, operators)
		if err != nil {
			return 0, fmt.Errorf("error checking if equation is possible: %w", err)
		}

		if possible {
			result += eq.result
		}
	}

	return result, nil
}
