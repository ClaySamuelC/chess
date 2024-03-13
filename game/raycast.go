package game

func CheckPiece(p *Piece, targetTeam string, targetRanks *map[string]bool) bool {
	return p.Team == targetTeam && (*targetRanks)[p.Rank] == true
}

func IsInBounds(pos int, dx int, dy int) bool {
	return pos%8+dx >= 0 && pos%8+dx < 8 && pos/8+dy >= 0 && pos/8+dy < 8
}

func (c *Chess) RayCast(pos int, d *Vector2, targetTeam string, targetRanks *map[string]bool) bool {
	// check bounds
	for step := 1; IsInBounds(pos, d.X*step, d.Y*step); step++ {
		p := c.Board[pos+d.Y*step*8+d.X*step]
		if p != nil {
			return CheckPiece(p, targetTeam, targetRanks)
		}
	}

	return false
}

func (c *Chess) CheckSquares(pos int, deltas *[]*Vector2, targetTeam string, targetRanks *map[string]bool) bool {
	for _, d := range *deltas {
		if IsInBounds(pos, d.X, d.Y) {
			p := c.Board[pos+d.Y*8+d.X]
			if p != nil {
				return p.Team == targetTeam && (*targetRanks)[p.Rank] == true
			}
		}
	}

	return false
}

func (c *Chess) MultiRayCast(pos int, deltas *[]*Vector2, targetTeam string, targetRanks *map[string]bool) bool {
	for _, d := range *deltas {
		if c.RayCast(pos, d, targetTeam, targetRanks) {
			return true
		}
	}
	return false
}
