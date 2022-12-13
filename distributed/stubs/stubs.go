package stubs

import "uk.ac.bris.cs/gameoflife/util"

var (
	GameOfLife = "Broker.GameOfLife"
	KeyPress   = "Broker.KeyPress"
	Next       = "Server.Next"
)

type WorldEntity struct {
	Width  int
	Height int
	Turns  int
	Grid   [][]bool
}

type GolRequest struct {
	World WorldEntity
}

type GolResponse struct {
	World WorldEntity
	Turns int
}

type KeyPressRequest struct {
	Key rune
}

type KeyPressResponse struct {
	World  WorldEntity
	Turn   int
	Paused bool
}

type NextRequest struct {
	World WorldEntity
	Cells []util.Cell
}

type NextResponse struct {
	Cells map[util.Cell]bool
}

type AliveCellsCountRequest struct {
}

type AliveCellsCountResponse struct {
	Count int
	Turn  int
}

type TurnCompleteRequest struct {
	Cells []util.Cell
	Turn  int
}

type TurnCompleteResponse struct{}
