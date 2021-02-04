package main

import (
	"fmt"

	"github.com/WeaselChess/Weasel/engine"
	"github.com/WeaselChess/Weasel/uci"
)

//EngineInfo is the info for our engine
var EngineInfo = uci.EngineInfo{
	Name:    "Weasel",
	Version: "v0.0.1-beta",
	Author:  "WeaselChess Club",
}

func init() {
	println("                                                  ")
	println("██╗    ██╗███████╗ █████╗ ███████╗███████╗██╗     ")
	println("██║    ██║██╔════╝██╔══██╗██╔════╝██╔════╝██║     ")
	println("██║ █╗ ██║█████╗  ███████║███████╗█████╗  ██║     ")
	println("██║███╗██║██╔══╝  ██╔══██║╚════██║██╔══╝  ██║     ")
	println("╚███╔███╔╝███████╗██║  ██║███████║███████╗███████╗")
	println(" ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝")
	println("                                                  ")
}

func main() {
	//uci.UCI(EngineInfo)
	//TODO: Make engine
	var playBitBoard uint64 = 0

	engine.SetBit(&playBitBoard, 61)
	engine.PrintBitBoard(playBitBoard)

	engine.ClearBit(&playBitBoard, 61)
	engine.PrintBitBoard(playBitBoard)

	var pos engine.BoardStruct
	err := pos.LoadFEN(engine.StartPosFEN)
	if err != nil {
		panic(err)
	}
	fmt.Println(pos)
}
