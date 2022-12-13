package gol

import (
	"sync"
	"uk.ac.bris.cs/gameoflife/util"
)

type World struct {
	WorldTurns   int
	WorldThreads int
	WorldWidth   int
	WorldHeight  int
	NowGrid      *WorldGrid
	PreviousGrid *WorldGrid
}

const OutRange = -1

// CellFliped send the flipped cell to the channel
func (world *World) CellFliped(turn int, channels distributorChannels) {
	for y := 0; y < world.WorldHeight; y++ {
		for x := 0; x < world.WorldWidth; x++ {
			world.cellCheck(x, y, turn, channels)
		}
	}
}

func (world *World) cellCheck(x, y, turn int, channels distributorChannels) {
	if world.WorldChange(x, y) {
		cell := util.Cell{x, y}
		cellFlipped := CellFlipped{
			CompletedTurns: turn,
			Cell:           cell,
		}
		channels.events <- cellFlipped
	}
}

// AliveCells AliveCount send the alive count to the channel
func (world *World) AliveCells(turn int, channel distributorChannels) {
	channel.events <- AliveCellsCount{
		CompletedTurns: turn,
		CellsCount:     world.GetAliveCount(),
	}
}

// Run runs the distributor.
func (world *World) Run() {
	wg := sync.WaitGroup{}
	for y := 0; y < world.WorldHeight; {
		for i := 0; i < world.WorldThreads; i++ {
			wg.Add(1)

			go world.waitGroupRun(&wg, y)

			y++
			if y == world.WorldHeight {
				break
			}
		}
	}
	wg.Wait()
	world.NowGrid, world.PreviousGrid = world.PreviousGrid, world.NowGrid
}

func (world *World) waitGroupRun(wg *sync.WaitGroup, y int) {
	defer wg.Done()
	for x := 0; x < world.WorldWidth; x++ {
		world.PreviousGrid.Set(x, y, world.NowGrid.NextState(x, y))
	}
}

// Set sets the state of the cell at the given coordinates.
func (g *WorldGrid) Set(x, y int, alive bool) {
	if g.CheckIn(x, y) {
		g.GridCells[x][y] = alive
	}
}

// CheckIn checks if the given coordinates are within the grid.
func (g *WorldGrid) CheckIn(x, y int) bool {
	return (x >= 0 && x < g.GridWidth && y >= 0 && y < g.GridHeight)
}

// CompletedTurns send the completed turn to the channel
func (world *World) CompletedTurns(turn int, channel distributorChannels) {
	channel.events <- TurnComplete{CompletedTurns: turn}
}

// FinalTurnComplete send the final turn to the channel
func (world *World) FinalTurnComplete(channel distributorChannels) {
	var cells []util.Cell
	for y := 0; y < world.WorldHeight; y++ {
		for x := 0; x < world.WorldWidth; x++ {
			if world.NowGrid.Get(x, y) {
				cells = append(cells, util.Cell{X: x, Y: y})
			}
		}
	}
	channel.events <- FinalTurnComplete{
		CompletedTurns: world.WorldTurns,
		Alive:          cells,
	}
}

// WorldChange IsChanged returns true if the cell at the given coordinates has changed.
func (world *World) WorldChange(x, y int) bool {
	return world.NowGrid.Get(x, y) != world.PreviousGrid.Get(x, y)
}

// Get returns the state of the cell at the given coordinates.
func (grid *WorldGrid) Get(coordX, coordY int) bool {
	if coordX >= 0 && coordX < grid.GridWidth && coordY >= 0 && coordY < grid.GridHeight {
		return grid.GridCells[coordX][coordY]
	} else {
		return grid.GetOutRange(coordX, coordY)
	}
}

func (grid *WorldGrid) GetOutRange(coordX, coordY int) bool {
	if coordX == OutRange || coordX == grid.GridWidth {
		coordX = (coordX + grid.GridWidth) % grid.GridWidth
	}
	if coordY == OutRange || coordY == grid.GridHeight {
		coordY = (coordY + grid.GridHeight) % grid.GridHeight
	}
	return grid.GridCells[coordX][coordY]
}

// Save saves the grid to the given file.
func (world *World) Save(filename string, channel distributorChannels) error {
	channel.ioCommand <- ioOutput
	channel.ioFilename <- filename
	for y := 0; y < world.WorldHeight; y++ {
		for x := 0; x < world.WorldHeight; x++ {
			if world.NowGrid.Get(x, y) {
				channel.ioOutput <- 255
			} else {
				channel.ioOutput <- 0
			}
		}
	}
	return nil
}

// GetAliveCount AliveCount returns the number of alive cells in the grid.
func (world *World) GetAliveCount() int {
	count := 0
	for y := 0; y < world.WorldHeight; y++ {
		for x := 0; x < world.WorldWidth; x++ {
			if world.NowGrid.Get(x, y) {
				count++
			}
		}
	}
	return count
}
