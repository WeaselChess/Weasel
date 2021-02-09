package uci

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/WeaselChess/Weasel/engine/board"
	"github.com/WeaselChess/Weasel/engine/search"
)

//EngineInfo holds the info for our engine
type EngineInfo struct {
	Name    string
	Version string
	Author  string
}

//Current board position
var pos board.PositionStruct

//Search Info
var info search.InfoStruct

//Set to true if a position is set
var positionSet bool

//UCI is our main loop for
func UCI(engineInfo EngineInfo) {
	var command []string
	var ready bool
	scanner := bufio.NewScanner(os.Stdin)

	space := regexp.MustCompile(`\s+`) //Used to delete multiple spaces

	for scanner.Scan() {
		index := 0

		command = strings.Split(space.ReplaceAllString(scanner.Text(), " "), " ")

	top:

		switch command[index] {
		case "uci":
			go uciHander(engineInfo)
		case "debug":
			go func() {
				if command[index+1] == "on" {
					board.DEBUG = true
				} else {
					board.DEBUG = false
				}
			}()
		case "isready":
			go func() {

				board.Initialize()
				//Init hash tables size with 2 MB's
				pos.PVTable.Init(2) //TODO: Add option to set hash size in mb

				ready = true

				fmt.Println("readyok")
			}()
		case "setoption":
		case "register":
		case "ucinewgame":
			pos.PVTable.Clear()
			err := pos.LoadFEN(board.StartPosFEN)
			if err != nil {
				panic(err)
			}
		case "position":
			if ready {
				go positionHandler(command[index:])
			}
		case "go":
			if ready {
				go goHandler(command[index:])
			}
		case "stop":
			info.Stopped = true
		case "ponderhit":
		case "quit":
			os.Exit(0)
		case "print":
			go pos.Print()
		case "divide":
			if ready && positionSet {
				go divideHander(command[index:])
			}
		default:
			if len(command) == index+1 {
				break
			}

			index++
			goto top
		}
	}
}

func uciHander(engineInfo EngineInfo) {
	fmt.Println("id name " + engineInfo.Name + engineInfo.Version)
	fmt.Println("id author " + engineInfo.Author)

	//TODO: Add supported options

	fmt.Println("uciok")
}

func positionHandler(command []string) {
	var boardSet bool
	if command[1] == "startpos" {
		err := pos.LoadFEN(board.StartPosFEN)
		if err != nil {
			panic(err)
		}
		boardSet = true
	} else if command[1] == "fen" {
		err := pos.LoadFEN(strings.Join(command[2:], " "))
		if err != nil {
			panic(err)
		}
		boardSet = true
	}

	str := strings.Join(command[2:], " ")

	if strings.Contains(str, "moves") && boardSet {

		index := strings.Index(str, "moves")

		str = str[index+6:]

		moves := strings.Split(str, " ")

		for i := 0; i < len(moves); i++ {
			move, err := pos.ParseMove(moves[i])
			if err != nil {
				panic(err)
			}

			if move != board.NoMove {
				var moveMade bool
				moveMade, err := pos.MakeMove(move)
				if err != nil {
					panic(err)
				}
				if moveMade {
					continue
				}
			}

			fmt.Printf("Non-legal move: %s", moves[i])

			break
		}

	}
	positionSet = true

}

func divideHander(command []string) {
	if len(command) > 1 {
		if unicode.IsDigit(rune(command[1][0])) {
			var depth int = int(rune(command[1][0]) - '0')
			err := pos.PerftDivide(depth)
			if err != nil {
				panic(err)
			}
		}
	}
}

func goHandler(command []string) {
	var err error
	depth := -1
	movesToGo := 30
	moveTime := -1
	timeV := -1
	inc := 0
	info.TimeSet = false
	for i := 0; i < len(command); i++ {
		switch command[i] {
		case "binc":
			if pos.Side == 1 {
				inc, err = strconv.Atoi(command[i+1])
				if err != nil {
					fmt.Println("Failed to parse binc time")
				}
			}
		case "winc":
			if pos.Side == 0 {
				inc, err = strconv.Atoi(command[i+1])
				if err != nil {
					fmt.Println("Failed to parse winc time")
				}
			}
		case "wtime":
			if pos.Side == 0 {
				timeV, err = strconv.Atoi(command[i+1])
				if err != nil {
					fmt.Println("Failed to parse wtime time")
				}
			}
		case "btime":
			if pos.Side == 1 {
				timeV, err = strconv.Atoi(command[i+1])
				if err != nil {
					fmt.Println("Failed to parse btime time")
				}
			}

		case "movestogo":
			movesToGo, err = strconv.Atoi(command[i+1])
			if err != nil {
				fmt.Println("Failed to parse binc time")
			}

		case "movetime":
			moveTime, err = strconv.Atoi(command[i+1])
			if err != nil {
				fmt.Println("Failed to parse binc time")
			}
		case "depth":
			depth, err = strconv.Atoi(command[i+1])
			if err != nil {
				fmt.Println("Failed to parse binc time")
			}
		}
	}

	if moveTime != -1 {
		timeV = moveTime
		movesToGo = 1
	}

	info.StartTime = time.Now().UnixNano() / int64(time.Millisecond)
	info.Depth = depth

	if timeV != -1 {
		info.TimeSet = true
		timeV /= movesToGo
		timeV -= 50
		info.StopTime = info.StartTime + int64(timeV+inc)
	}

	if depth == -1 {
		info.Depth = board.MaxDepth
	}

	fmt.Printf("time: %d start: %d stop: %d depth %d timeset %v\n",
		timeV, info.StartTime, info.StopTime, info.Depth, info.TimeSet)

	for {
		if positionSet {
			err := info.SearchPosition(&pos)
			if err != nil {
				panic(err)
			}
			positionSet = false
			break
		}
	}
}
