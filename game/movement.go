package game

import (
	"fmt"
)

func isValidMove(src int, moveX int, moveY int) bool {
	return !(src%8+moveX < 0 || src%8+moveX > 7 || src/8+moveY < 0 || src/8+moveY > 7)
}

func (c *Chess) canMoveInDirection(pos int, dirX int, dirY int, team string) bool {
	dest := pos + dirY*8 + dirX
	return isValidMove(pos, dirX, dirY) && (c.Board[dest] == nil || c.Board[dest].Team != team)
}

func (c *Chess) getMovesInDirection(pos int, dirX int, dirY int, team string) []int {
	moves := make([]int, 0)

	for true {
		if !c.canMoveInDirection(pos, dirX, dirY, team) {
			return moves
		}

		pos += dirY*8 + dirX
		moves = append(moves, pos)

		if c.Board[pos] != nil {
			return moves
		}
	}

	fmt.Println("Something's wrong")
	return moves
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

func (c *Chess) getKingMoves(pos int, team string) []int {
	moves := make([]int, 0)

	if c.canMoveInDirection(pos, -1, -1, team) {
		moves = append(moves, pos-9)
	}
	if c.canMoveInDirection(pos, -1, 0, team) {
		moves = append(moves, pos-1)
	}
	if c.canMoveInDirection(pos, -1, 1, team) {
		moves = append(moves, pos+7)
	}
	if c.canMoveInDirection(pos, 0, -1, team) {
		moves = append(moves, pos-8)
	}
	if c.canMoveInDirection(pos, 0, 1, team) {
		moves = append(moves, pos+8)
	}
	if c.canMoveInDirection(pos, 1, 0, team) {
		moves = append(moves, pos+1)
	}
	if c.canMoveInDirection(pos, 1, -1, team) {
		moves = append(moves, pos-7)
	}
	if c.canMoveInDirection(pos, 1, 1, team) {
		moves = append(moves, pos+9)
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
		moves = append(moves, pos+dy)
		if pos/8 == start && c.Board[pos+dy*2] == nil {
			moves = append(moves, pos+dy*2)
		}
	}

	// check forward-left
	if (pos%8-1 >= 0 && c.Board[pos+dy-1] != nil && c.Board[pos+dy-1].Team != team) || c.EnPassantLoc == pos+dy-1 {
		fmt.Println("Adding forward left.")
		moves = append(moves, pos+dy-1)
	}
	// check forward-right
	if (pos%8+1 <= 7 && c.Board[pos+dy+1] != nil && c.Board[pos+dy+1].Team != team) || c.EnPassantLoc == pos+dy+1 {
		moves = append(moves, pos+dy+1)
	}

	return moves
}

func (c *Chess) getRookMoves(pos int, team string) []int {
	moves := make([]int, 0)

	moves = append(moves, c.getMovesInDirection(pos, -1, 0, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 0, -1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 0, 1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 1, 0, team)...)

	return moves
}

func (c *Chess) getKnightMoves(pos int, team string) []int {
	moves := make([]int, 0)

	if c.canMoveInDirection(pos, -2, -1, team) {
		moves = append(moves, pos-10)
	}
	if c.canMoveInDirection(pos, -2, 1, team) {
		moves = append(moves, pos+6)
	}
	if c.canMoveInDirection(pos, -1, -2, team) {
		moves = append(moves, pos-17)
	}
	if c.canMoveInDirection(pos, -1, 2, team) {
		moves = append(moves, pos+15)
	}
	if c.canMoveInDirection(pos, 1, -2, team) {
		moves = append(moves, pos-15)
	}
	if c.canMoveInDirection(pos, 1, 2, team) {
		moves = append(moves, pos+17)
	}
	if c.canMoveInDirection(pos, 2, -1, team) {
		moves = append(moves, pos-6)
	}
	if c.canMoveInDirection(pos, 2, 1, team) {
		moves = append(moves, pos+10)
	}

	fmt.Printf("%v\n", moves)

	return moves
}

func (c *Chess) getBishopMoves(pos int, team string) []int {
	moves := make([]int, 0)

	moves = append(moves, c.getMovesInDirection(pos, -1, -1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, -1, 1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 1, -1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 1, 1, team)...)

	return moves
}

func (c *Chess) getQueenMoves(pos int, team string) []int {
	moves := make([]int, 0)

	moves = append(moves, c.getMovesInDirection(pos, -1, -1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, -1, 0, team)...)
	moves = append(moves, c.getMovesInDirection(pos, -1, 1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 0, -1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 0, 1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 1, 0, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 1, -1, team)...)
	moves = append(moves, c.getMovesInDirection(pos, 1, 1, team)...)

	return moves
}
