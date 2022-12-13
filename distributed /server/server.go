package main

import (
	"flag"
	"log"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/gameoflife/gol"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type Server struct {
}

func (s *Server) next(req stubs.NextRequest, res *stubs.NextResponse) (err error) {
	//receiver parameters
	params := gol.Params{
		ImageHeight: req.World.Height,
		ImageWidth:  req.World.Width,
		Turns:       req.World.Turns,
	}

	world := gol.World{WorldHeight: params.ImageHeight,
		WorldWidth: params.ImageWidth,
		WorldTurns: params.Turns,
	}

	world.NowGrid.GridCells = req.World.Grid

	res.Cells = make(map[util.Cell]bool)
	for _, cell := range req.Cells {
		res.Cells[cell] = world.NowGrid.NextState(cell.X, cell.Y)
	}
	return
}

func main() {
	//listen port on 8888
	network := "tcp"
	name := "port"
	port := "8888"
	usage := "listen port on"

	portNum := flag.String(name, port, usage)
	flag.Parse()
	listener, fail := net.Listen(network, ":"+*portNum)
	if fail == nil {
		//use defer release
		defer listener.Close()
		rpc.Register(new(Server))
		rpc.Accept(listener)
	} else {
		log.Fatalf("failed to listen: %v", fail)
	}
}
