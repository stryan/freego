package main

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
type Colour int

//Colour consts
const (
	Red Colour = iota
	Blue
	NoColour
)

//Game represents general game state
type Game struct {
	board *Board
	state GameState
}

//NewGame creates a new game and sets the state to lobby
func NewGame() *Game {
	return &Game{
		board: NewBoard(8),
		state: gameLobby,
	}
}

func (g *Game) String() string {
	var board string
	for i := 0; i < g.board.size; i++ {
		board = board + "|"
		for j := 0; j < g.board.size; j++ {
			board = board + g.board.board[i][j].String()
		}
		board = board + "|\n"
	}

	return fmt.Sprintf("Status: %v\n%v", g.state, board)
}

//Parse a given RawCommand,  return a ParsedCommand or an error
func (g *Game) Parse(p Player, cmd *RawCommand) (*ParsedCommand, error) {
	if g.state != gameTurnRed && g.state != gameTurnBlue {
		return nil, errors.New("game has not started")
	}
	if !g.board.validatePoint(cmd.srcX, cmd.srcY) || !g.board.validatePoint(cmd.dstX, cmd.dstY) {
		return nil, errors.New("invalid location in command")
	}
	if (p.Colour() == Red && g.state != gameTurnRed) || (p.Colour() == Blue && g.state != gameTurnBlue) {
		return nil, errors.New("not your turn")
	}
	start, err := g.board.GetPiece(cmd.srcX, cmd.srcY)
	if err != nil {
		return nil, err
	}
	if cmd.act != "-" && cmd.act != "x" {
		return nil, errors.New("invalid command action")
	}
	end, err := g.board.GetPiece(cmd.dstX, cmd.dstY)
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
		if g.state == gameTurnRed {
			g.state = gameTurnBlue
		} else if g.state == gameTurnBlue {
			g.state = gameTurnRed
		} else {
			return false, errors.New("bad game state")
		}
	}
	return true, nil
}

//SetupPiece adds a piece to the board during setup
func (g *Game) SetupPiece(x, y int, p *Piece) (bool, error) {
	if g.state != gameSetup {
		return false, errors.New("Trying to setup piece when  not in setup")
	}
	if !g.board.validatePoint(x, y) {
		return false, errors.New("Invalid location")
	}
	if p.Owner != g.board.GetColor(x, y) {
		return false, errors.New("Can't setup piece on enemy board")
	}
	return g.board.Place(x, y, p)
}

//Start start the game
func (g *Game) Start() {
	g.state = gameTurnRed
}

func (g *Game) move(x, y, s, t int) (bool, error) {
	startPiece, err := g.board.GetPiece(x, y)
	if err != nil {
		return false, err
	}
	endPiece, err := g.board.GetPiece(s, t)
	if err != nil {
		return false, err
	}
	if startPiece == nil {
		return false, errors.New("invalid start piece for move")
	}
	if endPiece != nil {
		return false, nil
	}
	res, err := g.board.Place(s, t, startPiece)
	if err != nil {
		return false, err
	}
	if res {
		err = g.board.Remove(x, y)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (g *Game) strike(x, y, s, t int) (bool, error) {
	startPiece, err := g.board.GetPiece(x, y)
	if err != nil {
		return false, err
	}
	endPiece, err := g.board.GetPiece(s, t)
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
		g.board.Remove(x, y)
	case 0:
		//tie
		g.board.Remove(x, y)
		g.board.Remove(s, t)
	case 1:
		//endPiece lost
		g.board.Remove(s, t)
		//scouts replace the piece that was destroyed
		if startPiece.Rank == 1 {
			g.board.Remove(x, y)
			g.board.Place(s, t, startPiece)
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
	if atk.Rank == 3 && def.Rank == 12 {
		return 1, nil
	}
	//spy hitting marshal
	if atk.Rank == 0 && def.Rank == 10 {
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
