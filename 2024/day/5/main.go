package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func main() {
	rules, updates, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	validSum := 0
	invalidSum := 0
	for _, update := range updates {
		if isValid, _ := isUpdateValid(rules, update); isValid {
			middleElement := update[len(update)/2]
			validSum += middleElement
		} else {
			fixedUpdate, err := fixUpdate(rules, update)
			if err != nil {
				fmt.Printf("error fixing update: %v", err)
				return
			}
			middleElement := fixedUpdate[len(fixedUpdate)/2]
			invalidSum += middleElement
		}
	}

	fmt.Printf("(Part one) Sum of middle elements of valid updates: %d\n", validSum)
	fmt.Printf("(Part two) Sum of middle elements of invalid updates: %d\n", invalidSum)
}

func readInput() (rules map[int][]int, updates [][]int, err error) {
	file, err := os.Open("./input.txt")
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rulesRe := regexp.MustCompile(`^\d+\|\d+$`)
	updatesRe := regexp.MustCompile(`^\d+(,\d+)*$`)
	rules = make(map[int][]int)
	updates = make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()

		if rulesRe.MatchString(line) {
			values := strings.Split(line, "|")

			value, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, nil, fmt.Errorf("error converting key to int: %w", err)
			}

			key, err := strconv.Atoi(values[1])
			if err != nil {
				return nil, nil, fmt.Errorf("error converting value to int: %w", err)
			}

			rules[key] = append(rules[key], value)
		}

		if updatesRe.MatchString(line) {
			values := strings.Split(line, ",")
			update := make([]int, 0)
			for _, value := range values {
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return nil, nil, fmt.Errorf("error converting update value to int: %w", err)
				}

				update = append(update, intValue)
			}
			updates = append(updates, update)
		}
	}

	return rules, updates, nil
}

func isUpdateValid(rules map[int][]int, update []int) (isValid bool, faultyIndex int) {
	reversedUpdate := make([]int, len(update))
	copy(reversedUpdate, update)
	slices.Reverse(reversedUpdate)

	for i, value := range reversedUpdate {
		numbersToCheck, ok := rules[value]
		if !ok {
			continue
		}

		remainingValues := reversedUpdate[:i]
		violatesRule := slices.ContainsFunc(remainingValues, func(i int) bool {
			return slices.Contains(numbersToCheck, i)
		})
		if violatesRule {
			return false, len(reversedUpdate) - i
		}
	}

	return true, -1
}

func fixUpdate(rules map[int][]int, update []int) ([]int, error) {
	updateCopy := make([]int, len(update))
	copy(updateCopy, update)

	if isValid, _ := isUpdateValid(rules, updateCopy); isValid {
		return updateCopy, nil
	}

	for {
		isValid, faultyIndex := isUpdateValid(rules, updateCopy)
		if isValid {
			return updateCopy, nil
		}
		tempCopy := make([]int, len(updateCopy))
		copy(tempCopy, updateCopy)

		if faultyIndex != 0 {
			faultyElement := updateCopy[faultyIndex]
			newUpdate := make([]int, 0)
			newUpdate = append(newUpdate, faultyElement)
			newUpdate = append(newUpdate, updateCopy[:faultyIndex]...)
			newUpdate = append(newUpdate, updateCopy[faultyIndex+1:]...)
			copy(updateCopy, newUpdate)
		}
		if slices.Equal(updateCopy, tempCopy) {
			return nil, fmt.Errorf("unable to fix update")
		}
	}
}
