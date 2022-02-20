package main

import "fmt"

func main() {
	//red := NewDummyPlayer(Red)
	//blue := NewDummyPlayer(Blue)
	g := NewGame()
	g.state = gameSetup
	printboardcolours(g)
	return
}

func addpiece(game *Game, rank int, c Colour, x int, y int) {
	res, err := game.SetupPiece(x, y, NewPieceFromInt(rank, c))
	if err != nil {
		panic(err)
	}
	if !res {
		panic("can't setup")
	}
}

func addriver(game *Game, x int, y int) {
	res, err := game.board.AddTerrain(x, y, 1)
	if err != nil {
		panic(err)
	}
	if !res {
		panic("can't river")
	}
}

func printboardcolours(g *Game) {
	for i := range g.board.board {
		for j := range g.board.board[i] {
			c := "X"
			if g.board.board[i][j].colour == Red {
				c = "R"
			} else if g.board.board[i][j].colour == Blue {
				c = "B"
			}
			fmt.Printf("%v", c)
		}
		fmt.Println("")
	}
}
