package game

func (c *Chess) RayCast(pos int, xDir int, yDir int, targetTeam string, targetRanks *map[string]bool) bool {
	// check bounds
	for step := 1; pos%8+xDir*step >= 0 && pos%8+xDir*step < 8 && pos/8+yDir*step >= 0 && pos/8+yDir*step < 8; step++ {
		p := c.Board[pos+yDir*step*8+xDir*step]
		if p != nil {
			return p.Team == targetTeam && (*targetRanks)[p.Rank] == true
		}
	}

	return false
}
