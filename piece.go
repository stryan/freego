package freego

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
	Bomb
	Unknown
)

//Piece :game piece
type Piece struct {
	Rank   Rank
	Owner  Colour
	Hidden bool
}

//NewPieceFromInt creates a new piece
func NewPieceFromInt(r int, o Colour) *Piece {
	return &Piece{
		Rank:   Rank(r),
		Owner:  o,
		Hidden: false,
	}
}

//NewPiece generates a piece by rank
func NewPiece(r Rank, o Colour) *Piece {
	return &Piece{
		Rank:   r,
		Owner:  o,
		Hidden: true,
	}
}

//NewUnknownPiece creates a new hidden piece
func NewUnknownPiece(o Colour) *Piece {
	return &Piece{
		Owner:  o,
		Rank:   Unknown,
		Hidden: true,
	}
}
