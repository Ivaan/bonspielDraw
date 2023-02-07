package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	gameIDIndex := 0
	drawNameIndex := 1
	sheetNameIndex := 2
	winnerToGameIDIndex := 3
	loserToGameIDIndex := 4

	readFile, err := os.Open("drawData.txt")
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	bonspiel := make([]game, 0)
	for fileScanner.Scan() {
		fields := strings.Fields(fileScanner.Text())
		gameID := fields[gameIDIndex]
		drawName := fields[drawNameIndex]
		sheetName := fields[sheetNameIndex]
		winnerToGameID := fields[winnerToGameIDIndex]
		loserToGameID := fields[loserToGameIDIndex]

		bonspiel = append(bonspiel, game{
			gameID:         gameID,
			drawName:       drawName,
			sheetName:      sheetName,
			winnerToGameID: winnerToGameID,
			loserToGameID:  loserToGameID,
		})
	}
	for _, g := range bonspiel {
		fmt.Println(g)
	}
}

type game struct {
	gameID         string
	drawName       string
	sheetName      string
	winnerToGameID string
	loserToGameID  string
}
