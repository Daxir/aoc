package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	stones, err := readInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	blinks := 25
	sum := 0
	for _, stone := range stones {
		sum += countStones(stone, blinks)
	}
	fmt.Printf("(Part one) Stone count: %d\n", sum)

	blinks = 75
	sum = 0
	for _, stone := range stones {
		sum += countStones(stone, blinks)
	}
	fmt.Printf("(Part two) Stone count: %d\n", sum)
}

func readInput() ([]int, error) {
	file, err := os.ReadFile("input.txt")
	if err != nil {
		return nil, err
	}
	numbers := strings.Split(string(file), " ")
	values := make([]int, len(numbers))
	for i, n := range numbers {
		value, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return values, nil
}

var cache = map[string]int{}

func countStones(stone, blinks int) int {
	key := fmt.Sprintf("s%d-b%d", stone, blinks)

	if hit, ok := cache[key]; ok {
		return hit
	}

	if blinks < 0 {
		return 0
	}

	if blinks == 0 {
		return 1
	}

	if stone == 0 {
		result := countStones(1, blinks-1)
		cache[key] = result
		return result
	}

	stoneAsString := strconv.Itoa(stone)
	if len(stoneAsString)%2 == 0 {
		left := stoneAsString[:len(stoneAsString)/2]
		right := stoneAsString[len(stoneAsString)/2:]
		leftAsInt, err := strconv.Atoi(left)
		if err != nil {
			panic(err)
		}
		rightAsInt, err := strconv.Atoi(right)
		if err != nil {
			panic(err)
		}
		result := countStones(leftAsInt, blinks-1) + countStones(rightAsInt, blinks-1)
		cache[key] = result
		return result
	}

	result := countStones(stone*2024, blinks-1)
	cache[key] = result
	return result
}
