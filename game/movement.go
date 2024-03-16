package game

import (
	"fmt"
)

func DeepCopyGame(c *Chess) *Chess {
	copy := &Chess{
		Board:         c.Board,
		Turn:          c.Turn,
		EnPassantLoc:  c.EnPassantLoc,
		HalfMoveClock: c.HalfMoveClock,
		FullMoveClock: c.FullMoveClock,
		PlayerInfo: map[string]*Info{
			"White": {c.PlayerInfo["White"].IsKingCastleValid, c.PlayerInfo["White"].IsQueenCastleValid, c.PlayerInfo["White"].KingPos},
			"Black": {c.PlayerInfo["Black"].IsKingCastleValid, c.PlayerInfo["Black"].IsQueenCastleValid, c.PlayerInfo["Black"].KingPos},
		},
	}

	return copy
}

func isValidMove(src int, move *Vector2) bool {
	return !(src%8+move.X < 0 || src%8+move.X > 7 || src/8+move.Y < 0 || src/8+move.Y > 7)
}

func (c *Chess) canMoveInDirection(pos int, d *Vector2, team string) bool {
	dest := pos + d.Y*8 + d.X
	return isValidMove(pos, d) && (c.Board[dest] == nil || c.Board[dest].Team != team)
}

func (c *Chess) getMovesInDirection(pos int, d *Vector2, team string) []int {
	moves := make([]int, 0)

	for {
		if !c.canMoveInDirection(pos, d, team) {
			return moves
		}

		pos += d.Y*8 + d.X
		moves = append(moves, pos)

		if c.Board[pos] != nil {
			return moves
		}
	}
}

func (c *Chess) GetPossibleMoves(p *Piece, pos int) []int {
	fmt.Printf("Checking for %v moves.\n", p.Rank)
	if p.Rank == "King" {
		return c.getKingMoves(pos, p.Team)
	}
	if p.Rank == "Pawn" {
		return c.getPawnMoves(pos, p.Team)
	}
	if p.Rank == "Rook" {
		return c.getRookMoves(pos, p.Team)
	}
	if p.Rank == "Knight" {
		return c.getKnightMoves(pos, p.Team)
	}
	if p.Rank == "Bishop" {
		return c.getBishopMoves(pos, p.Team)
	}
	if p.Rank == "Queen" {
		return c.getQueenMoves(pos, p.Team)
	}

	return nil
}

func (c *Chess) WillBeInCheck(src int, dest int, team string) bool {
	copy := DeepCopyGame(c)
	copy.Move(src, dest)

	kingPos := copy.PlayerInfo[team].KingPos

	fmt.Printf("Checking if player is in check at (%c%v)%v\n", 'a'+kingPos%8, 8-kingPos/8, kingPos)

	return copy.IsInCheck(kingPos, team)
}

func (c *Chess) getKingMoves(pos int, team string) []int {
	moves := make([]int, 0)

	for _, d := range *Adjacents {
		if c.canMoveInDirection(pos, d, team) {
			dest := pos + d.Y*8 + d.X

			if !c.WillBeInCheck(pos, dest, team) {
				moves = append(moves, dest)
			}
		}
	}

	if c.PlayerInfo[team].IsKingCastleValid {
		if c.canMoveInDirection(pos, Right, team) && c.canMoveInDirection(pos+1, Right, team) {
			if !c.WillBeInCheck(pos, pos+1, team) && !c.WillBeInCheck(pos, pos+2, team) {
				moves = append(moves, pos+2)
			}
		}
	}

	if c.PlayerInfo[c.Turn].IsQueenCastleValid {
		if c.canMoveInDirection(pos, Left, team) && c.canMoveInDirection(pos-1, Left, team) {
			if !c.WillBeInCheck(pos, pos-1, team) && !c.WillBeInCheck(pos, pos-2, team) {
				moves = append(moves, pos-2)
			}
		}
	}

	return moves
}

func (c *Chess) getPawnMoves(pos int, team string) []int {
	moves := make([]int, 0)

	start := 1
	dy := 8
	if team == "White" {
		dy = -8
		start = 6
	}

	if c.Board[pos+dy] == nil {
		if !c.WillBeInCheck(pos, pos+dy, team) {
			moves = append(moves, pos+dy)
		}
		if pos/8 == start && c.Board[pos+dy*2] == nil {
			if !c.WillBeInCheck(pos, pos+dy*2, team) {
				moves = append(moves, pos+dy*2)
			}
		}
	}

	// check forward-left
	if (pos%8-1 >= 0 && c.Board[pos+dy-1] != nil && c.Board[pos+dy-1].Team != team) || c.EnPassantLoc == pos+dy-1 {
		if !c.WillBeInCheck(pos, pos+dy-1, team) {
			moves = append(moves, pos+dy-1)
		}
	}
	// check forward-right
	if (pos%8+1 <= 7 && c.Board[pos+dy+1] != nil && c.Board[pos+dy+1].Team != team) || c.EnPassantLoc == pos+dy+1 {
		if !c.WillBeInCheck(pos, pos+dy+1, team) {
			moves = append(moves, pos+dy+1)
		}
	}

	return moves
}

func (c *Chess) getRookMoves(pos int, team string) []int {
	moves := make([]int, 0)

	for _, d := range *Straights {
		for _, move := range c.getMovesInDirection(pos, d, team) {
			if !c.WillBeInCheck(pos, move, team) {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (c *Chess) getKnightMoves(pos int, team string) []int {
	moves := make([]int, 0)

	for _, d := range *KnightMoves {
		if c.canMoveInDirection(pos, d, team) {
			if !c.WillBeInCheck(pos, pos+d.Y*8+d.X, team) {
				moves = append(moves, pos+d.Y*8+d.X)
			}
		}
	}

	return moves
}

func (c *Chess) getBishopMoves(pos int, team string) []int {
	moves := make([]int, 0)

	for _, d := range *Diagonals {
		for _, move := range c.getMovesInDirection(pos, d, team) {
			if !c.WillBeInCheck(pos, move, team) {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (c *Chess) getQueenMoves(pos int, team string) []int {
	moves := make([]int, 0)

	for _, d := range *Adjacents {
		for _, move := range c.getMovesInDirection(pos, d, team) {
			if !c.WillBeInCheck(pos, move, team) {
				moves = append(moves, move)
			}
		}
	}

	return moves
}
