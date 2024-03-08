package game

type Piece struct {
	Team string
	Rank string
}

func NewPiece(rank string, team string) *Piece {
	return &Piece{team, rank}
}
