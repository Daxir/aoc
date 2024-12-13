package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const tenBillion int = 10000000000000

func main() {
	configs, err := readInput()
	if err != nil {
		fmt.Println(err)
	}

	totalCost := 0
	for _, config := range configs {
		winner, err := getPressesToPrize(config)
		if err != nil {
			continue
		}
		totalCost += getTokenCost(winner)
	}

	fmt.Printf("(Part one) Total cost: %d\n", totalCost)

	totalCost = 0
	for _, config := range configs {
		newConfig := machineConfig{config.a, config.b, coordinate{config.prize.x + tenBillion, config.prize.y + tenBillion}}
		winner, err := getPressesToPrize(newConfig)
		if err != nil {
			continue
		}
		totalCost += getTokenCost(winner)
	}
	fmt.Printf("(Part two) Total cost: %d\n", totalCost)
}

type coordinate struct {
	x, y int
}

type machineConfig struct {
	a, b, prize coordinate
}

func readInput() ([]machineConfig, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	onEmptyLine := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == '\n' {
				if i+1 < len(data) && data[i+1] == '\n' {
					return i + 2, data[:i], nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		return 0, data, bufio.ErrFinalToken
	}
	scanner.Split(onEmptyLine)

	machineConfigs := make([]machineConfig, 0)
	for scanner.Scan() {
		configuration := scanner.Text()
		aRe := regexp.MustCompile(`Button A:.*\n`)
		bRe := regexp.MustCompile(`Button B:.*\n`)
		prizeRe := regexp.MustCompile(`Prize:.*`)
		a := aRe.FindString(configuration)
		b := bRe.FindString(configuration)
		prize := prizeRe.FindString(configuration)

		a = strings.ReplaceAll(a, "Button A: ", "")
		a = strings.Trim(a, "\n")
		a, _ = strings.CutPrefix(a, "X+")
		aCoords := strings.Split(a, ", Y+")
		ax, err := strconv.Atoi(aCoords[0])
		if err != nil {
			return nil, err
		}
		ay, err := strconv.Atoi(aCoords[1])
		if err != nil {
			return nil, err
		}
		aButton := coordinate{ax, ay}

		b = strings.ReplaceAll(b, "Button B: ", "")
		b = strings.Trim(b, "\n")
		b, _ = strings.CutPrefix(b, "X+")
		bCoords := strings.Split(b, ", Y+")
		bx, err := strconv.Atoi(bCoords[0])
		if err != nil {
			return nil, err
		}
		by, err := strconv.Atoi(bCoords[1])
		if err != nil {
			return nil, err
		}
		bButton := coordinate{bx, by}

		prize = strings.ReplaceAll(prize, "Prize: ", "")
		prize = strings.Trim(prize, "\n")
		prize, _ = strings.CutPrefix(prize, "X=")
		prizeCoords := strings.Split(prize, ", Y=")
		px, err := strconv.Atoi(prizeCoords[0])
		if err != nil {
			return nil, err
		}
		py, err := strconv.Atoi(prizeCoords[1])
		if err != nil {
			return nil, err
		}
		prizeLocation := coordinate{px, py}

		machineConfigs = append(machineConfigs, machineConfig{aButton, bButton, prizeLocation})
	}

	return machineConfigs, nil
}

func getTokenCost(winner coordinate) int {
	return winner.x*3 + winner.y
}

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

func getY(config machineConfig) float64 {
	by := float64(config.b.y)
	ax := float64(config.a.x)
	ay := float64(config.a.y)
	bx := float64(config.b.x)
	py := float64(config.prize.y)
	px := float64(config.prize.x)

	y := (ay*px - ax*py) / (bx*ay - ax*by)

	if !isIntegral(y) {
		return math.NaN()
	}
	return y
}

func getX(config machineConfig, y float64) float64 {
	ax := float64(config.a.x)
	bx := float64(config.b.x)
	px := float64(config.prize.x)

	x := (px - bx*y) / ax

	if !isIntegral(x) {
		return math.NaN()
	}

	return x
}

func getPressesToPrize(config machineConfig) (coordinate, error) {
	y := getY(config)
	if math.IsNaN(y) {
		return coordinate{}, fmt.Errorf("y is not an integer")
	}

	x := getX(config, y)
	if math.IsNaN(x) {
		return coordinate{}, fmt.Errorf("x is not an integer")
	}

	return coordinate{int(x), int(y)}, nil
}
