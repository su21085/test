package gol

import (
	"fmt"
	"time"
)

// QUIT Define constants to indicate meaning
const QUIT = 'q'
const SAVE = 's'
const PAUSED = 'p'

// Define more channel queue
type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	keyPresses <-chan rune
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(param Params, channel distributorChannels) {
	// TODO: Create a 2D slice to store the world.
	var world = woldConfig(param, channel)

	turn := 0
	world.CellFliped(turn, channel) // send the initial state to the channel
	// TODO: Execute all turns of the Game of Life.
	trunTicker := time.NewTicker(time.Second * 2)
	for turn < param.Turns {
		select {
		case <-trunTicker.C:
			world.AliveCells(turn, channel)
		case keyPress := <-channel.keyPresses:
			switch keyPress {
			case SAVE:
				filename := fmt.Sprintf("%dx%dx%d-%d", param.ImageHeight, param.ImageWidth, param.Threads, turn)
				world.Save(filename, channel)
				image := ImageOutputComplete{turn, filename}
				channel.events <- image
			case PAUSED:
				stage := StateChange{CompletedTurns: turn, NewState: Paused}
				channel.events <- stage
				for {
					if <-channel.keyPresses == PAUSED {
						channel.events <- StateChange{turn, Executing}
						break
					}
				}
			case QUIT:
				stage := StateChange{CompletedTurns: turn, NewState: Quitting}
				channel.ioCommand <- ioCheckIdle
				<-channel.ioIdle
				channel.events <- stage
				close(channel.events)
				return

			}
		default:
			turn = turn + 1
			world.Run()
			world.CellFliped(turn, channel)
			world.CompletedTurns(turn, channel)
		}
	}
	// TODO: Report the final state using FinalTurnCompleteEvent.
	world.FinalTurnComplete(channel)
	filename := fmt.Sprintf("%dx%dx%d-%d", param.ImageHeight, param.ImageWidth, param.Threads, turn)
	world.Save(filename, channel)
	channel.events <- ImageOutputComplete{
		Filename:       filename,
		CompletedTurns: turn,
	}

	// Make sure that the Io has finished any output before exiting.
	channel.ioCommand <- ioCheckIdle
	<-channel.ioIdle

	channel.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(channel.events)
}

//create new world to store 2D slice

// init world on game start
func woldConfig(params Params, disChannel distributorChannels) *World {
	var world World
	world.WorldWidth = params.ImageWidth
	world.WorldHeight = params.ImageHeight
	world.WorldTurns = params.Turns
	world.WorldThreads = params.Threads
	world.NowGrid = initGrid(params)
	world.PreviousGrid = initGrid(params)

	//print log to console
	fmt.Println("wold init successfully,world:", world)

	disChannel.ioCommand <- ioInput // ioInput 1
	disChannel.ioFilename <- fmt.Sprintf("%dx%d", params.ImageHeight, params.ImageWidth)
	for y := 0; y < world.WorldHeight; y++ {
		for x := 0; x < world.WorldWidth; x++ {
			intChat := <-disChannel.ioInput
			if x >= 0 && x < world.NowGrid.GridWidth && y >= 0 && y < world.NowGrid.GridHeight {
				world.NowGrid.GridCells[x][y] = intChat == 255
			}
		}
	}

	return &world
}

// init graid on game start by given params
func initGrid(params Params) *WorldGrid {
	var worldGrid WorldGrid
	worldGrid.GridHeight = params.ImageHeight
	worldGrid.GridWidth = params.ImageWidth
	worldGrid.GridThreads = params.Threads
	worldGrid.GridCells = make([][]bool, worldGrid.GridHeight)
	for index := range worldGrid.GridCells {
		worldGrid.GridCells[index] = make([]bool, worldGrid.GridWidth)
	}
	return &worldGrid
}

func reportTheFinalState(world *World, c distributorChannels, params Params, turn int) {
	world.FinalTurnComplete(c)
	filename := fmt.Sprintf("%dx%dx%d-%d", params.ImageHeight, params.ImageWidth, params.Threads, turn)
	world.Save(filename, c)
	c.events <- ImageOutputComplete{
		Filename:       filename,
		CompletedTurns: turn,
	}
}

func turnsOfTheGameOfLife(turn int, param Params, world *World, channel distributorChannels) {
	trunTicker := time.NewTicker(time.Second * 2)
	for turn < param.Turns {
		select {
		case <-trunTicker.C:
			world.AliveCells(turn, channel)
		case keyPress := <-channel.keyPresses:
			switch keyPress {
			case SAVE:
				filename := fmt.Sprintf("%dx%dx%d-%d", param.ImageHeight, param.ImageWidth, param.Threads, turn)
				world.Save(filename, channel)
				image := ImageOutputComplete{turn, filename}
				channel.events <- image
			case PAUSED:
				stage := StateChange{CompletedTurns: turn, NewState: Paused}
				channel.events <- stage
				for {
					if <-channel.keyPresses == PAUSED {
						channel.events <- StateChange{turn, Executing}
						break
					}
				}
			case QUIT:
				stage := StateChange{CompletedTurns: turn, NewState: Quitting}
				channel.ioCommand <- ioCheckIdle
				<-channel.ioIdle
				channel.events <- stage
				close(channel.events)
				return

			}
		default:
			turn = turn + 1
			world.Run()
			world.CellFliped(turn, channel)
			world.CompletedTurns(turn, channel)
		}
	}
}

func CellFliped(turn int, world *World, channel distributorChannels) {
	turn = turn + 1
	world.Run()
	world.CellFliped(turn, channel)
	world.CompletedTurns(turn, channel)
}
