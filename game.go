package freego

import (
	"errors"
	"fmt"
	"math"
)

//GameState is an int representing the current state the game's in
type GameState int

const (
	gameLobby GameState = iota
	gameSetup
	gameTurnRed
	gameTurnBlue
	gameWinRed
	gameWinBlue
)

//Colour reperensents player colour
//go:generate enumer -type=Colour
type Colour int

//Colour consts
const (
	Red Colour = iota
	Blue
	NoColour
)

//Game represents general game state
type Game struct {
	Board *Board
	State GameState
}

//NewGame creates a new game and sets the state to lobby
func NewGame() *Game {
	return &Game{
		Board: NewBoard(8),
		State: gameLobby,
	}
}

func (g *Game) String() string {
	var board string
	for i := 0; i < g.Board.size; i++ {
		board = board + "|"
		for j := 0; j < g.Board.size; j++ {
			board = board + g.Board.board[i][j].String()
		}
		board = board + "|\n"
	}

	return fmt.Sprintf("Status: %v\n%v", g.State, board)
}

//Parse a given RawCommand,  return a ParsedCommand or an error
func (g *Game) Parse(p Player, cmd *RawCommand) (*ParsedCommand, error) {
	if g.State != gameTurnRed && g.State != gameTurnBlue {
		return nil, errors.New("game has not started")
	}
	if !g.Board.validatePoint(cmd.srcX, cmd.srcY) || !g.Board.validatePoint(cmd.dstX, cmd.dstY) {
		return nil, errors.New("invalid location in command")
	}
	if (p.Colour() == Red && g.State != gameTurnRed) || (p.Colour() == Blue && g.State != gameTurnBlue) {
		return nil, errors.New("not your turn")
	}
	start, err := g.Board.GetPiece(cmd.srcX, cmd.srcY)
	if err != nil {
		return nil, err
	}
	if cmd.act != "-" && cmd.act != "x" {
		return nil, errors.New("invalid command action")
	}
	end, err := g.Board.GetPiece(cmd.dstX, cmd.dstY)
	if err != nil {
		return nil, err
	}
	if start == nil {
		return nil, errors.New("empty start piece")
	}
	startRank := start.Rank.String()
	endRank := ""
	if end != nil {
		endRank = end.Rank.String()
	}
	return &ParsedCommand{cmd.srcX, cmd.srcY, startRank, cmd.dstX, cmd.dstY, endRank, cmd.act}, nil
}

//Mutate the game state given a ParsedCommand
func (g *Game) Mutate(cmd *ParsedCommand) (bool, error) {
	var res bool
	var err error
	switch cmd.act {
	case "-":
		res, err = g.move(cmd.srcX, cmd.srcY, cmd.dstX, cmd.dstY)
	case "x":
		res, err = g.strike(cmd.srcX, cmd.srcY, cmd.dstX, cmd.dstY)
	default:
		return false, errors.New("invalid mutate action")
	}
	if err != nil {
		return false, err
	}
	if res {
		if g.State == gameTurnRed {
			g.State = gameTurnBlue
		} else if g.State == gameTurnBlue {
			g.State = gameTurnRed
		} else {
			return false, errors.New("bad game state")
		}
	}
	return res, nil
}

//SetupPiece adds a piece to the board during setup
func (g *Game) SetupPiece(x, y int, p *Piece) (bool, error) {
	if g.State != gameSetup {
		return false, errors.New("Trying to setup piece when  not in setup")
	}
	if p == nil {
		return false, errors.New("Tried to setup a nil piece")
	}
	if !g.Board.validatePoint(x, y) {
		return false, errors.New("Invalid location")
	}
	if p.Owner != g.Board.GetColor(x, y) {
		return false, fmt.Errorf("Can't setup piece on enemy board: %v != %v", p.Owner, g.Board.GetColor(x, y))
	}
	return g.Board.Place(x, y, p)
}

//Start start the game
func (g *Game) Start() bool {
	if g.State == gameSetup {
		g.State = gameTurnRed
		return true
	}
	return false
}

//Setup puts the game in setup mode
func (g *Game) Setup() bool {
	if g.State == gameLobby {
		g.State = gameSetup
		return true
	}
	return false
}

func (g *Game) move(x, y, s, t int) (bool, error) {
	startPiece, err := g.Board.GetPiece(x, y)
	if err != nil {
		return false, err
	}
	endPiece, err := g.Board.GetPiece(s, t)
	if err != nil {
		return false, err
	}
	if startPiece == nil {
		return false, errors.New("invalid start piece for move")
	}
	if endPiece != nil {
		return false, nil
	}
	//attempt to remove starting piece first
	err = g.Board.Remove(x, y)
	if err != nil {
		return false, err
	}
	// then place piece in new location
	res, err := g.Board.Place(s, t, startPiece)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (g *Game) strike(x, y, s, t int) (bool, error) {
	startPiece, err := g.Board.GetPiece(x, y)
	if err != nil {
		return false, err
	}
	endPiece, err := g.Board.GetPiece(s, t)
	if err != nil {
		return false, err
	}
	if startPiece == nil {
		return false, errors.New("invalid start piece for strike")
	}
	if endPiece == nil {
		return false, nil
	}

	distance := math.Sqrt(math.Pow(float64(s-x), 2) + math.Pow(float64(t-y), 2))
	if distance > 1 && startPiece.Rank != 1 {
		//trying to attack a piece more then 1 square away when not a scout
		return false, nil
	}
	//bombs can't attack
	if startPiece.Rank == 12 {
		return false, nil
	}

	r, err := g.combat(startPiece, endPiece)
	if err != nil {
		return false, err
	}
	switch r {
	case -1:
		//startPiece lost
		err = g.Board.Remove(x, y)
		if err != nil {
			return true, err
		}
	case 0:
		//tie
		err = g.Board.Remove(x, y)
		err2 := g.Board.Remove(s, t)
		if err != nil || err2 != nil {
			return true, fmt.Errorf("Errors: %v %v", err, err2)
		}
	case 1:
		//endPiece lost
		err := g.Board.Remove(s, t)
		if err != nil {
			return true, err
		}
		//scouts replace the piece that was destroyed
		if startPiece.Rank == Scout {
			err = g.Board.Remove(x, y)
			if err != nil {
				return true, err
			}
			res, err := g.Board.Place(s, t, startPiece)
			if err != nil {
				return true, err
			}
			if !res {
				return false, errors.New("Combat was valid but somehow placing the new piece was not")
			}
		}
	}
	return true, nil
}

func (g *Game) combat(atk, def *Piece) (int, error) {
	if atk == nil || def == nil {
		return 0, errors.New("invalid attacker or defender")
	}
	if atk.Hidden {
		return 0, errors.New("trying to attack with a piece we know nothing about?")
	}
	if def.Hidden {
		return 0, errors.New("defender has not been revealed to us")
	}
	//handle special cases first
	//miner hitting bomb
	if atk.Rank == Miner && def.Rank == Bomb {
		return 1, nil
	}
	//anyone else hitting bomb
	if def.Rank == Bomb {
		return -1, nil
	}
	//spy hitting marshal
	if atk.Rank == Spy && def.Rank == Marshal {
		return 1, nil
	}
	//normal cases
	if atk.Rank > def.Rank {
		return 1, nil
	} else if atk.Rank == def.Rank {
		return 0, nil
	} else {
		return -1, nil
	}
}
