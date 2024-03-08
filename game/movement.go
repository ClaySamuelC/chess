package game

import (
	"fmt"
)

func (c *Chess) GetPossibleMoves(p *Piece, pos int) []int {
	fmt.Printf("Checking for %v moves.\n", p.Rank)
	// if p.Rank == "King" {
	// 	return c.getKingMoves(pos, p.Team)
	// }
	// if p.Rank == "Pawn" {
	// 	return c.getPawnMoves(pos, p.Team)
	// }
	// if p.Rank == "Rook" {
	// 	return c.getRookMoves(pos, p.Team)
	// }
	// if p.Rank == "Knight" {
	// 	return c.getKnightMoves(pos, p.Team)
	// }
	// if p.Rank == "Bishop" {
	// 	return c.getBishopMoves(pos, p.Team)
	// }
	// if p.Rank == "Queen" {
	// 	return c.getQueenMoves(pos, p.Team)
	// }

	return nil
}

func (c *Chess) getKingMoves(pos int, team string) []int {
	moves := make([]int, 0)

	return moves
}

func (c *Chess) getPawnMoves(pos int, team string) []int {
	moves := make([]int, 0)

	return moves
}

func (c *Chess) getRookMoves(pos int, team string) []int {
	moves := make([]int, 0)

	return moves
}

func (c *Chess) getKnightMoves(pos int, team string) []int {
	moves := make([]int, 0)

	return moves
}

func (c *Chess) getBishopMoves(pos int, team string) []int {
	moves := make([]int, 0)

	return moves
}

func (c *Chess) getQueenMoves(pos int, team string) []int {
	moves := make([]int, 0)

	return moves
}
