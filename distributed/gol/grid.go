package gol

// WorldGrid create cell by Grid
type WorldGrid struct {
	GridWidth   int
	GridHeight  int
	GridThreads int
	GridCells   [][]bool
}

// CheckIn checks if the given coordinates are within the grid.
func (grid *WorldGrid) withinTheGrid(x, y int) bool {
	return x >= 0 && x < grid.GridWidth && y >= 0 && y < grid.GridHeight
}

func (grid *WorldGrid) NextState(x, y int) bool {
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if grid.Get(x+i, y+j) && (i != 0 || j != 0) {
				alive++
			}
		}
	}
	return alive == 3 || (alive == 2 && grid.Get(x, y))
}
