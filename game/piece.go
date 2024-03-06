package game

type Piece struct {
	Rank string
	Team string
}

func NewPiece(rank string, team string) *Piece {
	return &Piece{rank, team}
}

func (p *Piece) ToString() string {
	var str string = string(p.Team[0])
	if p.Rank == "Knight" {
		str += "N"
	} else {
		str += string(p.Rank[0])
	}
	return str
}
