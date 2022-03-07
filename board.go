package freego

import "errors"

//Board represents main game board
type Board struct {
	board [][]*Tile
	size  int
}

//NewBoard creates a new board instance
func NewBoard(size int) *Board {
	if size < 4 || size%2 != 0 {
		return nil
	}
	b := make([][]*Tile, size)
	var colour Colour
	for i := 0; i < size; i++ {
		b[i] = make([]*Tile, size)
		if i < size/2 {
			colour = Blue
		} else {
			colour = Red
		}
		for j := 0; j < size; j++ {
			b[i][j] = &Tile{i, j, true, nil, colour}
		}
	}
	return &Board{
		board: b,
		size:  size,
	}
}

func (b *Board) validatePoint(x, y int) bool {
	if x < 0 || x >= len(b.board) {
		return false
	}
	if y < 0 || y >= len(b.board) {
		return false
	}
	return true
}

//GetPiece returns the piece at a given location
func (b *Board) GetPiece(x, y int) (*Piece, error) {
	if !b.validatePoint(x, y) {
		return nil, errors.New("GetPiece invalid location")
	}
	return b.board[y][x].Piece(), nil
}

//Place a piece on the board; returns false if a piece is already there
func (b *Board) Place(x, y int, p *Piece) (bool, error) {
	if !b.validatePoint(x, y) {
		return false, errors.New("Place invalid location")
	}
	if b.board[y][x].Piece() != nil {
		return false, nil
	}
	err := b.board[y][x].Place(p)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Remove a piece from a tile
func (b *Board) Remove(x, y int) error {
	if !b.validatePoint(x, y) {
		return errors.New("Remove invalid location")
	}
	b.board[y][x].Remove()
	return nil
}

//GetColor returns color of tile
func (b *Board) GetColor(x, y int) Colour {
	return b.board[y][x].Colour()
}

//AddTerrain puts a river tile at specified location
func (b *Board) AddTerrain(x, y, t int) (bool, error) {
	if !b.validatePoint(x, y) {
		return false, errors.New("River invalid location")
	}
	b.board[y][x].AddTerrain(t)
	return true, nil
}

//IsTerrain checks if tile is terrain
func (b *Board) IsTerrain(x, y int) (bool, error) {
	if !b.validatePoint(x, y) {
		return false, errors.New("River check invalid location")
	}
	return !b.board[y][x].passable, nil
}
