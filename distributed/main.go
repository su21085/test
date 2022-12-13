package main

import (
	"flag"
	"fmt"
	"runtime"

	"uk.ac.bris.cs/gameoflife/gol"
	"uk.ac.bris.cs/gameoflife/sdl"
)

// main is the function called when starting Game of Life with 'go run .'
func main() {
	runtime.LockOSThread()
	var params gol.Params

	flag.IntVar(
		&params.Threads,
		"t",
		8,
		"Specify the number of worker threads to use. Defaults to 8.")

	flag.IntVar(
		&params.ImageWidth,
		"w",
		512,
		"Specify the width of the image. Defaults to 512.")

	flag.IntVar(
		&params.ImageHeight,
		"h",
		512,
		"Specify the height of the image. Defaults to 512.")

	flag.IntVar(
		&params.Turns,
		"turns",
		10000000000,
		"Specify the number of turns to process. Defaults to 10000000000.")

	noVis := flag.Bool(
		"noVis",
		false,
		"Disables the SDL window, so there is no visualisation during the tests.")

	flag.Parse()

	fmt.Println("Threads:", params.Threads)
	fmt.Println("Width:", params.ImageWidth)
	fmt.Println("Height:", params.ImageHeight)

	//协程之间的通讯
	keyPresses := make(chan rune, 10)
	events := make(chan gol.Event, 1000)
	//入口函数，因为所有的程序都有一个入口函数，用来启动成的函数。
	//并发进行
	go gol.Run(params, events, keyPresses) //键盘字符出队列
	if !(*noVis) {
		//前端图形库 去接收键盘字符 'p' 's' 'q' 'k'
		sdl.Run(params, events, keyPresses) //键盘字符入队
	} else {
		complete := false
		for !complete {
			event := <-events
			switch event.(type) {
			case gol.FinalTurnComplete:
				complete = true
			}
		}
	}
}
