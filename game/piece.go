package game

type Piece struct {
	Rank string
	Team string
}

func NewPiece(rank string, team string) *Piece {
	return &Piece{rank, team}
}
