package game

import (
	"fmt"
)

func (b *Board) isValidDest(x int, y int, team string) bool {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		fmt.Println("False")
		return false
	}

	destPiece := b.Squares[y][x]
	if destPiece != nil && destPiece.Team == team {
		return false
	}

	return true
}

func (b *Board) appendIfValid(xMove int, yMove int, arr []*Vector2, pos *Vector2, team string) ([]*Vector2, string) {
	if team == "White" {
		yMove = -yMove
	}

	dest := &Vector2{pos.X + xMove, pos.Y + yMove}

	if b.isValidDest(dest.X, dest.Y, team) {
		if b.Squares[dest.Y][dest.X] == nil {
			return append(arr, dest), "Blank"
		}

		return append(arr, dest), b.Squares[dest.Y][dest.X].Team
	}

	return arr, "Invalid"
}

func (b *Board) GetPossibleMoves(p *Piece, pos *Vector2) []*Vector2 {
	fmt.Printf("Checking for %v moves.\n", p.Rank)
	if p.Rank == "King" {
		return b.getKingMoves(pos, p.Team)
	}
	if p.Rank == "Pawn" {
		return b.getPawnMoves(pos, p.Team)
	}
	if p.Rank == "Rook" {
		return b.getRookMoves(pos, p.Team)
	}
	if p.Rank == "Knight" {
		return b.getKnightMoves(pos, p.Team)
	}
	if p.Rank == "Bishop" {
		return b.getBishopMoves(pos, p.Team)
	}
	if p.Rank == "Queen" {
		return b.getQueenMoves(pos, p.Team)
	}

	return nil
}

func (b *Board) getKingMoves(pos *Vector2, team string) []*Vector2 {
	moves := make([]*Vector2, 0)

	moves, _ = b.appendIfValid(-1, -1, moves, pos, team)
	moves, _ = b.appendIfValid(-1, 0, moves, pos, team)
	moves, _ = b.appendIfValid(-1, 1, moves, pos, team)
	moves, _ = b.appendIfValid(0, -1, moves, pos, team)
	moves, _ = b.appendIfValid(0, 1, moves, pos, team)
	moves, _ = b.appendIfValid(1, -1, moves, pos, team)
	moves, _ = b.appendIfValid(1, 0, moves, pos, team)
	moves, _ = b.appendIfValid(1, 1, moves, pos, team)

	return moves
}

func (b *Board) getPawnMoves(pos *Vector2, team string) []*Vector2 {
	moves := make([]*Vector2, 0)
	var moveSignal string

	_, moveSignal = b.appendIfValid(0, 1, moves, pos, team)
	if moveSignal == "Blank" {
		moves, _ = b.appendIfValid(0, 1, moves, pos, team)

		if moveSignal == "Blank" && pos.Y == 6 && team == "White" || pos.Y == 1 && team == "Black" {
			moves, _ = b.appendIfValid(0, 2, moves, pos, team)
		}
	}

	_, moveSignal = b.appendIfValid(-1, 1, moves, pos, team)
	if moveSignal != "Blank" && moveSignal != team {
		moves, _ = b.appendIfValid(-1, 1, moves, pos, team)
		// if en passant is valid
	} else if b.CanEnPassant(pos, -1, 1) {
		fmt.Println("Can En Passant!")
		moves, _ = b.appendIfValid(-1, 1, moves, pos, team)
	}
	_, moveSignal = b.appendIfValid(1, 1, moves, pos, team)
	if moveSignal != "Blank" && moveSignal != team {
		moves, _ = b.appendIfValid(1, 1, moves, pos, team)
	} else if b.CanEnPassant(pos, 1, 1) {
		fmt.Println("Can En Passant!")
		moves, _ = b.appendIfValid(1, 1, moves, pos, team)
	}

	return moves
}

func (b *Board) getRookMoves(pos *Vector2, team string) []*Vector2 {
	moves := make([]*Vector2, 0)
	var moveSignal string

	directions := [4]Vector2{
		{0, 1},
		{0, -1},
		{-1, 0},
		{1, 0},
	}

	for _, dir := range directions {
		flag := true

		for step := 1; flag; step++ {
			moves, moveSignal = b.appendIfValid(step*dir.X, step*dir.Y, moves, pos, team)
			flag = moveSignal == "Blank"
		}
	}

	return moves
}

func (b *Board) getKnightMoves(pos *Vector2, team string) []*Vector2 {
	moves := make([]*Vector2, 0)

	moves, _ = b.appendIfValid(2, 1, moves, pos, team)
	moves, _ = b.appendIfValid(2, -1, moves, pos, team)
	moves, _ = b.appendIfValid(1, 2, moves, pos, team)
	moves, _ = b.appendIfValid(1, -2, moves, pos, team)
	moves, _ = b.appendIfValid(-1, 2, moves, pos, team)
	moves, _ = b.appendIfValid(-1, -2, moves, pos, team)
	moves, _ = b.appendIfValid(-2, 1, moves, pos, team)
	moves, _ = b.appendIfValid(-2, -1, moves, pos, team)

	return moves
}

func (b *Board) getBishopMoves(pos *Vector2, team string) []*Vector2 {
	moves := make([]*Vector2, 0)
	var moveSignal string

	directions := [4]Vector2{
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, dir := range directions {
		flag := true

		for step := 1; flag; step++ {
			moves, moveSignal = b.appendIfValid(step*dir.X, step*dir.Y, moves, pos, team)
			flag = moveSignal == "Blank"
		}
	}

	return moves
}

func (b *Board) getQueenMoves(pos *Vector2, team string) []*Vector2 {
	moves := make([]*Vector2, 0)
	var moveSignal string

	directions := [8]Vector2{
		{0, 1},
		{0, -1},
		{-1, 0},
		{1, 0},
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, dir := range directions {
		flag := true

		for step := 1; flag; step++ {
			moves, moveSignal = b.appendIfValid(step*dir.X, step*dir.Y, moves, pos, team)
			flag = moveSignal == "Blank"
		}
	}

	return moves
}
