package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	robots, err := readInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	movementTime := 100
	xBound := 101
	yBound := 103

	for j := 0; j < len(robots); j++ {
		robots[j].move(movementTime)
		robots[j].wrap(xBound, yBound)
	}

	safetyMap := constructSafetyMap(robots, xBound, yBound)
	safetyFactor := safetyMap.getSafetyFactor()

	fmt.Printf("(Part one) Safety factor: %d\n", safetyFactor)

	for j := 0; j < len(robots); j++ {
		robots[j].reset()
	}

	verticalLineTime := findVerticalLineTime(robots, xBound, yBound)
	if verticalLineTime == -1 {
		fmt.Println("No vertical line found")
	} else {
		fmt.Printf("(Part two) Vertical line found at time: %d\n", verticalLineTime)
	}
}

type coordinate struct {
	x int
	y int
}

type robot struct {
	startingPosition coordinate
	velocity         coordinate
	currentPosition  coordinate
}

func (r *robot) move(seconds int) {
	r.currentPosition.x += r.velocity.x * seconds
	r.currentPosition.y += r.velocity.y * seconds
}

func (r *robot) wrap(xBound, yBound int) {
	r.currentPosition.x = ((r.currentPosition.x % xBound) + xBound) % xBound
	r.currentPosition.y = ((r.currentPosition.y % yBound) + yBound) % yBound
}

func (r *robot) reset() {
	r.currentPosition = r.startingPosition
}

func readInput() ([]robot, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var robots []robot
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		vectors := strings.Split(line, " ")
		positionVector, _ := strings.CutPrefix(vectors[0], "p=")
		velocityVector, _ := strings.CutPrefix(vectors[1], "v=")
		positionValues := strings.Split(positionVector, ",")
		velocityValues := strings.Split(velocityVector, ",")
		px, err := strconv.Atoi(positionValues[0])
		if err != nil {
			return nil, err
		}
		py, err := strconv.Atoi(positionValues[1])
		if err != nil {
			return nil, err
		}
		vx, err := strconv.Atoi(velocityValues[0])
		if err != nil {
			return nil, err
		}
		vy, err := strconv.Atoi(velocityValues[1])
		if err != nil {
			return nil, err
		}

		robots = append(robots, robot{
			startingPosition: coordinate{x: px, y: py},
			velocity:         coordinate{x: vx, y: vy},
			currentPosition:  coordinate{x: px, y: py},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return robots, nil
}

type safetyMap struct {
	q1, q2, q3, q4 []*robot
}

func constructSafetyMap(robots []robot, xBound, yBound int) safetyMap {
	// split the map into 4 quadrants, elements exactly on the boundary are ignored
	q1 := make([]*robot, 0)
	q2 := make([]*robot, 0)
	q3 := make([]*robot, 0)
	q4 := make([]*robot, 0)

	midX := xBound / 2
	midY := yBound / 2

	for i := 0; i < len(robots); i++ {
		r := &robots[i]

		// ignore elements exactly on the boundary (only for even bounds)
		if xBound%2 == 0 && r.currentPosition.x == midX {
			continue
		}
		switch {
		case r.currentPosition.x < midX && r.currentPosition.y < midY:
			q1 = append(q1, r)
		case r.currentPosition.x > midX && r.currentPosition.y < midY:
			q2 = append(q2, r)
		case r.currentPosition.x < midX && r.currentPosition.y > midY:
			q3 = append(q3, r)
		case r.currentPosition.x > midX && r.currentPosition.y > midY:
			q4 = append(q4, r)
		}
	}

	return safetyMap{
		q1: q1,
		q2: q2,
		q3: q3,
		q4: q4,
	}
}

func (s *safetyMap) getSafetyFactor() int {
	return len(s.q1) * len(s.q2) * len(s.q3) * len(s.q4)
}

func drawRobots(robots []robot, xBound, yBound int) {
	grid := make([][]rune, yBound)
	for i := 0; i < yBound; i++ {
		grid[i] = make([]rune, xBound)
		for j := 0; j < xBound; j++ {
			grid[i][j] = '.'
		}
	}

	for i := 0; i < len(robots); i++ {
		r := &robots[i]
		grid[r.currentPosition.y][r.currentPosition.x] = '#'
	}

	for i := 0; i < yBound; i++ {
		for j := 0; j < xBound; j++ {
			fmt.Printf("%c", grid[i][j])
		}
		fmt.Println()
	}
}

func findVerticalLineTime(robots []robot, xBound, yBound int) int {
	maxTime := math.MaxInt32
	for i := range maxTime {
		for j := 0; j < len(robots); j++ {
			robots[j].move(1)
			robots[j].wrap(xBound, yBound)
		}

		if isVerticalLine(robots) {
			drawRobots(robots, xBound, yBound)
			return i + 1
		}
	}
	return -1
}

func isVerticalLine(robots []robot) bool {
	minLineSize := 10
	columns := make(map[int][]int)

	for _, r := range robots {
		columns[r.currentPosition.x] = append(columns[r.currentPosition.x], r.currentPosition.y)
	}

	for _, yCoords := range columns {
		slices.Sort(yCoords)

		count := 1
		for i := 1; i < len(yCoords); i++ {
			if yCoords[i]-yCoords[i-1] <= 1 {
				count++
				if count >= minLineSize {
					return true
				}
			} else {
				count = 1
			}
		}
	}

	return false
}
