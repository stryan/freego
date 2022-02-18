package main

import "strconv"

//Rank represents the rank of a piece
type Rank int

func (r Rank) String() string {
	return strconv.Itoa(int(r))
}

//Piece :game piece
type Piece struct {
	Rank   Rank
	Owner  Colour
	Hidden bool
}

//NewPiece creates a new piece
func NewPiece(r int, o Colour) *Piece {
	return &Piece{
		Rank:   Rank(r),
		Owner:  o,
		Hidden: false,
	}
}

//NewHiddenPiece creates a new hidden piece
func NewHiddenPiece(o Colour) *Piece {
	return &Piece{
		Owner:  o,
		Hidden: true,
	}
}
