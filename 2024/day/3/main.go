package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	instructions, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	sum := 0
	for _, instruction := range instructions {
		result, err := executeInstruction(instruction)
		if err != nil {
			fmt.Printf("error executing instruction: %v", err)
			return
		}
		if result.instruction == mul {
			sum += result.value
		}
	}

	fmt.Printf("(Part one) Sum: %v\n", sum)

	sum = 0
	isEnabled := true
	for _, instruction := range instructions {
		result, err := executeInstruction(instruction)
		if err != nil {
			fmt.Printf("error executing instruction: %v", err)
			return
		}
		switch result.instruction {
		case do:
			isEnabled = true
		case dont:
			isEnabled = false
		case mul:
			if isEnabled {
				sum += result.value
			}
		default:
			fmt.Printf("invalid instruction: %v", result.instruction)
			return
		}
	}

	fmt.Printf("(Part two) Sum: %v\n", sum)
}

func readInput() ([]string, error) {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	re := regexp.MustCompile(`(mul\(\d{1,3},\d{1,3}\))|(do\(\))|(don't\(\))`)
	matches := re.FindAllString(string(file), -1)

	return matches, nil
}

type instruction string

var (
	mul instruction = "mul"
	do instruction = "do"
	dont instruction = "don't"
)

type instructionOutput struct {
	instruction instruction
	value			 int
}

func executeInstruction(instruction string) (instructionOutput, error) {
	instructionRegex := regexp.MustCompile(`[a-z']+`)
	instructionName := instructionRegex.FindString(instruction)
	
	switch instructionName {
	case "mul":
		argumentRegex := regexp.MustCompile(`\d{1,3}`)
		arguments := argumentRegex.FindAllString(instruction, -1)
		if len(arguments) != 2 {
			return instructionOutput{}, fmt.Errorf("invalid number of arguments: %v", arguments)
		}

		left, err := strconv.Atoi(arguments[0])
		if err != nil {
			return instructionOutput{}, fmt.Errorf("error converting to int: %v", err)
		}
		right, err := strconv.Atoi(arguments[1])
		if err != nil {
			return instructionOutput{}, fmt.Errorf("error converting to int: %v", err)
		}
		return instructionOutput{instruction: mul, value: left * right}, nil
	case "do":
		return instructionOutput{instruction: do}, nil
	case "don't":
		return instructionOutput{instruction: dont}, nil
	default:
		return instructionOutput{}, fmt.Errorf("invalid instruction: %v", instructionName)
	}
}