package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"slices"
)

func main() {
	diskSpace, err := readInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	compressed, err := compress(diskSpace)
	if err != nil {
		fmt.Println(err)
		return
	}

	checksum := calculateChecksum(compressed)
	fmt.Printf("(Part one) Compressed disk checksum: %v\n", checksum)

	chunks := splitIntoChunks(diskSpace)

	compressedChunks, err := compressChunks(chunks)
	if err != nil {
		fmt.Println(err)
		return
	}

	checksum = calculateChecksum(compressedChunks)
	fmt.Printf("(Part two) Compressed disk checksum: %v\n", checksum)
}

var digits = map[string]int{
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
}

func readInput() ([]int, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	diskSpace := make([]int, 0)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for id := 0; true; id++ {
		if !scanner.Scan() {
			break
		}
		fileBlock := scanner.Text()
		digit, ok := digits[fileBlock]
		if !ok {
			return nil, fmt.Errorf("invalid digit: %s", fileBlock)
		}
		for i := 0; i < digit; i++ {
			diskSpace = append(diskSpace, id)
		}

		if !scanner.Scan() {
			break
		}
		freeSpace := scanner.Text()
		digit, ok = digits[freeSpace]
		if !ok {
			return nil, fmt.Errorf("invalid digit: %s", freeSpace)
		}
		for i := 0; i < digit; i++ {
			diskSpace = append(diskSpace, -1)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return diskSpace, nil
}

func compress(diskSpace []int) ([]int, error) {
	compressed := make([]int, 0)
	next, stop := iter.Pull2(slices.Backward(diskSpace))
	defer stop()

	j := len(diskSpace) - 1
	for i, cell := range diskSpace {
		if i >= j {
			break
		}

		if cell != -1 {
			compressed = append(compressed, cell)
			continue
		}

		for {
			k, nextCell, ok := next()
			j = k
			if !ok {
				return nil, fmt.Errorf("unexpected end of disk space")
			}
			if nextCell != -1 {
				compressed = append(compressed, nextCell)
				break
			}
		}
	}

	return compressed, nil
}

func calculateChecksum(diskSpace []int) int {
	sum := 0
	for pos, block := range diskSpace {
		if block == -1 {
			continue
		}
		sum += pos * block
	}

	return sum
}

type chunk struct {
	fileID int
	size   int
}

func (c chunk) isFree() bool {
	return c.fileID == -1
}

func splitIntoChunks(diskSpace []int) []chunk {
	chunks := make([]chunk, 0)
	for i, block := range diskSpace {
		if i == 0 || block != diskSpace[i-1] {
			chunks = append(chunks, chunk{fileID: block, size: 1})
			continue
		}

		chunks[len(chunks)-1].size++
	}

	return chunks
}

func compressChunks(chunks []chunk) ([]int, error) {
	compressed := make([]chunk, 0)
	compressed = append(compressed, chunks...)

	for _, fileChunk := range slices.Backward(chunks) {
		if fileChunk.isFree() {
			continue
		}

		currentChunkIndex := slices.Index(compressed, fileChunk)
		if currentChunkIndex == -1 {
			return nil, fmt.Errorf("chunk not found: %v", fileChunk)
		}

		firstSuitableChunkIndex := slices.IndexFunc(compressed, func(c chunk) bool {
			return c.isFree() && c.size >= fileChunk.size
		})
		if firstSuitableChunkIndex == -1 || firstSuitableChunkIndex >= currentChunkIndex {
			continue
		}

		compressed[currentChunkIndex] = chunk{fileID: -1, size: fileChunk.size}

		remainingSize := compressed[firstSuitableChunkIndex].size - fileChunk.size
		compressed[firstSuitableChunkIndex] = fileChunk
		if remainingSize > 0 {
			newFreeChunk := chunk{fileID: -1, size: remainingSize}
			compressed = slices.Insert(compressed, firstSuitableChunkIndex+1, newFreeChunk)
		}
	}

	diskSpace := make([]int, 0)
	for _, chunk := range compressed {
		for i := 0; i < chunk.size; i++ {
			diskSpace = append(diskSpace, chunk.fileID)
		}
	}

	return diskSpace, nil
}
