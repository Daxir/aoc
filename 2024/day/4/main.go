package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input, err := readInput()
	if err != nil {
		fmt.Printf("error reading input: %v", err)
		return
	}

	substring := "XMAS"
	count := findSubstringInAllDirections(input, substring, settings{
		shouldSearchHorizontally: true,
		shouldSearchVertically:   true,
		shouldSearchDiagonally:   true,
	})
	fmt.Printf("(Part one) Substring %v appears %v times\n", substring, count)

	substring = "MAS"
	submatrices := createAllSquareSubmatrices(input, len(substring))
	count = 0
	for _, submatrix := range submatrices {
		subCount := findSubstringInAllDirections(submatrix, substring, settings{
			shouldSearchHorizontally: false,
			shouldSearchVertically:   false,
			shouldSearchDiagonally:   true,
		})
		// If the substring appears in the diagonal twice, it means that it creates a "cross" with the center
		if subCount == 2 {
			count++
		}
	}

	fmt.Printf("(Part two) A cross of the substring %v appears %v times\n", substring, count)
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

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	height := len(text)
	width := len(text[0])
	for i, row := range text {
		if len(row) != width {
			return nil, fmt.Errorf("row %v is not the same length as the first row", i)
		}
	}
	if height == 0 || width == 0 {
		return nil, fmt.Errorf("no data found in file")
	}
	if height != width {
		return nil, fmt.Errorf("height and width are not the same")
	}

	return text, nil
}

type settings struct {
	shouldSearchHorizontally bool
	shouldSearchVertically   bool
	shouldSearchDiagonally   bool
}

func findSubstringInAllDirections(text [][]rune, substring string, settings settings) int {
	letters := []rune(substring)
	reversedLetters := make([]rune, len(letters))
	for i := 0; i < len(letters); i++ {
		reversedLetters[i] = letters[len(letters)-i-1]
	}
	reversedSubstring := string(reversedLetters)

	size := len(text)

	count := 0

	if settings.shouldSearchHorizontally {
		for _, row := range text {
			count += findInChunk(row, target{forwards: substring, backwards: reversedSubstring})
		}
	}

	if settings.shouldSearchVertically {
		for i := 0; i < size; i++ {
			column := make([]rune, size)
			for j := 0; j < size; j++ {
				column[j] = text[j][i]
			}
			count += findInChunk(column, target{forwards: substring, backwards: reversedSubstring})
		}
	}

	if settings.shouldSearchDiagonally {
		for i := 0; i < size; i++ {
			leftToRightUpperDiagonal := make([]rune, size)
			leftToRightLowerDiagonal := make([]rune, size)
			rightToLeftUpperDiagonal := make([]rune, size)
			rightToLeftLowerDiagonal := make([]rune, size)
			for j := 0; j < size-i; j++ {
				leftToRightUpperDiagonal[j] = text[j][j+i]
				leftToRightLowerDiagonal[j] = text[j+i][j]
				rightToLeftUpperDiagonal[j] = text[j][size-j-i-1]
				rightToLeftLowerDiagonal[j] = text[j+i][size-j-1]
			}
			count += findInChunk(leftToRightUpperDiagonal, target{forwards: substring, backwards: reversedSubstring})
			count += findInChunk(rightToLeftUpperDiagonal, target{forwards: substring, backwards: reversedSubstring})
			if i != 0 {
				count += findInChunk(leftToRightLowerDiagonal, target{forwards: substring, backwards: reversedSubstring})
				count += findInChunk(rightToLeftLowerDiagonal, target{forwards: substring, backwards: reversedSubstring})
			}
		}
	}

	return count
}

type target struct {
	forwards  string
	backwards string
}

func findInChunk(chunk []rune, target target) (count int) {
	count = 0
	for i := 0; i < len(chunk)-len(target.forwards)+1; i++ {
		chunk := string(chunk[i : i+len(target.forwards)])
		if chunk == target.forwards || chunk == target.backwards {
			count++
		}
	}
	return count
}

func createAllSquareSubmatrices(text [][]rune, size int) [][][]rune {
	submatrices := make([][][]rune, 0)
	for i := 0; i < len(text)-size+1; i++ {
		for j := 0; j < len(text[0])-size+1; j++ {
			submatrix := make([][]rune, size)
			for k := 0; k < size; k++ {
				submatrix[k] = make([]rune, size)
				for l := 0; l < size; l++ {
					submatrix[k][l] = text[i+k][j+l]
				}
			}
			submatrices = append(submatrices, submatrix)
		}
	}
	return submatrices
}
