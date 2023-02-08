package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	bonspiel := loadBonspiel("drawData.txt")
	teams := loadTeams("teamsData.txt")

	// for k, g := range bonspiel {
	// 	fmt.Println(k, g)
	// }

	// for _, t := range teams {
	// 	fmt.Println(t)
	// }
	var computIf func(g *game) []string
	computIf = func(g *game) []string {
		out := make([]string, 0)
		if g.winnerTo != nil {
			wins := computIf(g.winnerTo)
			for _, w := range wins {
				out = append(out, g.gameID+" w "+w)
			}
		} else {
			out = append(out, g.gameID+" w "+g.winnerToGameID)
		}

		if g.loserTo != nil {
			lossess := computIf(g.loserTo)
			for _, l := range lossess {
				out = append(out, g.gameID+" l "+l)
			}
		} else {
			out = append(out, g.gameID+" l "+g.loserToGameID)
		}
		return out
	}

	var getPaths func(g *game) [][]*game
	getPaths = func(g *game) [][]*game {
		out := make([][]*game, 0)
		if g.winnerTo != nil {
			wins := getPaths(g.winnerTo)
			for _, w := range wins {
				p := make([]*game, 0)
				p = append(p, g)
				p = append(p, w...)
				out = append(out, p)
			}
		} else {
			p := make([]*game, 0)
			p = append(p, g)
			out = append(out, p)
		}
		if g.loserTo != nil {
			losses := getPaths(g.loserTo)
			for _, l := range losses {
				q := make([]*game, 0)
				q = append(q, g)
				q = append(q, l...)
				out = append(out, q)
			}
		} else {
			q := make([]*game, 0)
			q = append(q, g)
			out = append(out, q)
		}

		return out
	}

	// team1computIf := computIf(bonspiel[teams[0].startingGameID])
	// for _, m := range team1computIf {
	// 	fmt.Println(m)
	// }
	// fmt.Println()
	// team1Paths := getPaths(bonspiel[teams[0].startingGameID])
	// for _, gs := range team1Paths {
	// 	for _, g := range gs {
	// 		fmt.Print(g.drawName, g.sheetName, " ")
	// 	}
	// 	fmt.Println()
	// }
	// teamPaths := make([][]*game, 0)
	// for _, t := range teams {
	// 	paths := getPaths(bonspiel[t.startingGameID])
	// 	teamPaths = append(teamPaths, paths...)
	// }

	// for _, gs := range teamPaths {
	// 	for _, g := range gs {
	// 		//fmt.Print(g.drawName, g.sheetName, " ")
	// 		fmt.Print(g.sheetName[4:], " ")
	// 	}
	// 	fmt.Println()
	// }

	var printTeamTree func(g *game, wl string, tabs int)
	printTeamTree = func(g *game, wl string, tabs int) {
		for i := 0; i < tabs; i++ {
			fmt.Print("\t")
		}
		fmt.Println(wl, "-", g.drawName)
		for i := 0; i < tabs; i++ {
			fmt.Print("\t")
		}
		fmt.Println(" ", " ", g.sheetName)

		if g.winnerTo != nil {
			printTeamTree(g.winnerTo, "w", tabs+1)
		} else {
			for i := 0; i < tabs+1; i++ {
				fmt.Print("\t")
			}
			fmt.Println(g.winnerToGameID)
		}

		if g.loserTo != nil {
			printTeamTree(g.loserTo, "l", tabs+1)
		} else {
			for i := 0; i < tabs+1; i++ {
				fmt.Print("\t")
			}
			fmt.Println(g.loserToGameID)
		}
	}
	printTeamTree(bonspiel[teams[0].startingGameID], "s", 0)
}

func loadBonspiel(drawDataFileName string) map[string]*game {
	gameIDIndex := 0
	drawNameIndex := 1
	sheetNameIndex := 2
	winnerToGameIDIndex := 3
	loserToGameIDIndex := 4

	readFile, err := os.Open(drawDataFileName)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	games := make([]game, 0)
	for fileScanner.Scan() {
		fields := strings.Fields(fileScanner.Text())
		gameID := fields[gameIDIndex]
		drawName := fields[drawNameIndex]
		sheetName := fields[sheetNameIndex]
		winnerToGameID := fields[winnerToGameIDIndex]
		loserToGameID := fields[loserToGameIDIndex]

		games = append(games, game{
			gameID:         gameID,
			drawName:       drawName,
			sheetName:      sheetName,
			winnerToGameID: winnerToGameID,
			loserToGameID:  loserToGameID,
		})
	}
	bonspiel := make(map[string]*game)
	for i, g := range games {
		bonspiel[g.gameID] = &games[i]
	}
	for k := range bonspiel {
		g := bonspiel[k]
		g.winnerTo = bonspiel[g.winnerToGameID]
		g.loserTo = bonspiel[g.loserToGameID]
		//fmt.Println("tada?", g)
	}
	// for k, g := range bonspiel {
	// 	fmt.Println(k, g)
	// }
	return bonspiel
}
func loadTeams(teamsDataFileName string) []team {
	nameIndex := 0
	startingGameIDIndex := 1

	readFile, err := os.Open(teamsDataFileName)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	teams := make([]team, 0)
	for fileScanner.Scan() {
		fields := strings.Fields(fileScanner.Text())
		name := fields[nameIndex]
		startingGameID := fields[startingGameIDIndex]

		teams = append(teams, team{
			name:           name,
			startingGameID: startingGameID,
		})
	}
	return teams
}

type game struct {
	gameID         string
	drawName       string
	sheetName      string
	winnerToGameID string
	loserToGameID  string
	winnerTo       *game
	loserTo        *game
}

type team struct {
	name           string
	startingGameID string
}

type ifWinGame struct {
	win    bool
	gameID string
}
