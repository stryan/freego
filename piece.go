package main

//Rank represents the rank of a piece
//go:generate enumer -type=Rank
type Rank int

//Rank rank pieces
const (
	Flag Rank = iota
	Spy
	Scout
	Miner
	Captain
	General
	Marshal
)

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
