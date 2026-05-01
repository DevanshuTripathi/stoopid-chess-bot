package main

import (
	"math/rand"
	"sort"
	"time"
)

const (
	CHECKMATE = 10000
	STALEMATE = 0
	// We can easily bump this to 4 in Go without Move Ordering,
	// but let's leave it at 3 just to see how instantaneous it feels compared to Python!
	DEPTH = 5
)

// Maps in Go are great for simple lookups
var pieceScore = map[byte]int{
	'K': 0, 'k': 0,
	'Q': 90, 'q': 90,
	'R': 50, 'r': 50,
	'B': 30, 'b': 30,
	'N': 30, 'n': 30,
	'P': 10, 'p': 10,
}

// --- PIECE SQUARE TABLES ---

var knightScores = [8][8]int{
	{-5, -4, -3, -3, -3, -3, -4, -5},
	{-4, -2, 0, 0, 0, 0, -2, -4},
	{-3, 0, 1, 1, 1, 1, 0, -3},
	{-3, 0, 1, 2, 2, 1, 0, -3},
	{-3, 0, 1, 2, 2, 1, 0, -3},
	{-3, 0, 1, 1, 1, 1, 0, -3},
	{-4, -2, 0, 0, 0, 0, -2, -4},
	{-5, -4, -3, -3, -3, -3, -4, -5},
}

var bishopScores = [8][8]int{
	{-2, -1, -1, -1, -1, -1, -1, -2},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-1, 0, 0, 1, 1, 0, 0, -1},
	{-1, 0, 1, 1, 1, 1, 0, -1},
	{-1, 0, 1, 1, 1, 1, 0, -1},
	{-1, 0, 0, 1, 1, 0, 0, -1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-2, -1, -1, -1, -1, -1, -1, -2},
}

var rookScores = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{1, 2, 2, 2, 2, 2, 2, 1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{0, 0, 0, 1, 1, 0, 0, 0},
}

var queenScores = [8][8]int{
	{-2, -1, -1, -1, -1, -1, -1, -2},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-1, 0, 1, 1, 1, 1, 0, -1},
	{-1, 0, 1, 1, 1, 1, 0, -1},
	{-1, 0, 1, 1, 1, 1, 0, -1},
	{-1, 0, 0, 1, 1, 0, 0, -1},
	{-1, 0, 0, 0, 0, 0, 0, -1},
	{-2, -1, -1, -1, -1, -1, -1, -2},
}

var pawnScores = [8][8]int{
	{8, 8, 8, 8, 8, 8, 8, 8},
	{8, 8, 8, 8, 8, 8, 8, 8},
	{5, 6, 6, 7, 7, 6, 6, 5},
	{2, 3, 3, 5, 5, 3, 3, 2},
	{1, 2, 3, 4, 4, 3, 2, 1},
	{1, 1, 2, 3, 3, 2, 1, 1},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

func init() {
	// Seed the random number generator on startup
	rand.Seed(time.Now().UnixNano())
}

// --- MINIMAX ALGORITHM ---

func FindBestMove(gs *GameState, validMoves []Move) Move {
	var bestMove Move
	turnMultiplier := 1
	if !gs.WhiteToMove {
		turnMultiplier = -1
	}

	orderMoves(gs, validMoves)
	negamax(gs, DEPTH, -CHECKMATE, CHECKMATE, turnMultiplier, &bestMove)
	return bestMove
}

func negamax(gs *GameState, depth int, alpha, beta int, turnMultiplier int, bestMove *Move) int {
	if depth == 0 {
		// Quiescence search prevents the bot from making "blind" trades
		return quiescenceSearch(gs, alpha, beta, turnMultiplier)
	}

	moves := gs.GetValidMoves()
	if len(moves) == 0 {
		if gs.Checkmate {
			return -(CHECKMATE + depth) // Current player is mated
		}
		return STALEMATE
	}
	orderMoves(gs, moves)

	maxScore := -CHECKMATE
	for _, move := range moves {
		gs.MakeMove(move)
		// Standard Negamax: negative result of recursive call with swapped alpha/beta
		score := -negamax(gs, depth-1, -beta, -alpha, -turnMultiplier, bestMove)
		gs.UndoMove()

		if score > maxScore {
			maxScore = score
			if depth == DEPTH {
				*bestMove = move
			}
		}

		if score > alpha {
			alpha = score
		}

		if alpha >= beta {
			break // Pruning
		}
	}
	return maxScore
}

func FindMoveMinimax(gs *GameState, validMoves []Move, depth int, alpha int, beta int, whiteToMove bool, bestMove *Move) int {
	if depth == 0 {
		return ScoreBoard(gs)
	}

	if whiteToMove {
		maxScore := -CHECKMATE
		for _, move := range validMoves {
			gs.MakeMove(move)
			nextMoves := gs.GetValidMoves()
			score := FindMoveMinimax(gs, nextMoves, depth-1, alpha, beta, false, bestMove)
			gs.UndoMove()

			if score > maxScore {
				maxScore = score
				if depth == DEPTH {
					*bestMove = move // Update the pointer with the best move
				}
			}
			alpha = max(alpha, maxScore)
			if beta <= alpha {
				break // Prune!
			}
		}
		return maxScore
	} else {
		minScore := CHECKMATE
		for _, move := range validMoves {
			gs.MakeMove(move)
			nextMoves := gs.GetValidMoves()
			score := FindMoveMinimax(gs, nextMoves, depth-1, alpha, beta, true, bestMove)
			gs.UndoMove()

			if score < minScore {
				minScore = score
				if depth == DEPTH {
					*bestMove = move
				}
			}
			beta = min(beta, minScore)
			if beta <= alpha {
				break // Prune!
			}
		}
		return minScore
	}
}

// --- BOARD EVALUATION ---

func ScoreBoard(gs *GameState) int {
	if gs.Stalemate {
		return STALEMATE
	}

	score := 0
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := gs.Board[r][c]
			if piece == '-' {
				continue
			}

			val := pieceScore[piece]
			posScore := 0
			switch piece {
			case 'N':
				posScore = knightScores[r][c]
			case 'n':
				posScore = knightScores[7-r][c]
			case 'B':
				posScore = bishopScores[r][c]
			case 'b':
				posScore = bishopScores[7-r][c]
			case 'R':
				posScore = rookScores[r][c]
			case 'r':
				posScore = rookScores[7-r][c]
			case 'Q':
				posScore = queenScores[r][c]
			case 'q':
				posScore = queenScores[7-r][c]
			case 'P':
				posScore = pawnScores[r][c]
			case 'p':
				posScore = pawnScores[7-r][c]
			}

			if isWhite(piece) {
				score += val*10 + posScore
			} else {
				score -= val*10 + posScore
			}
		}
	}
	return score
}

func orderMoves(gs *GameState, moves []Move) {
	for i := range moves {
		score := 0
		if moves[i].PieceCaptured != '-' {
			// Prioritize captures of valuable pieces by less valuable ones
			score = 10*pieceScore[moves[i].PieceCaptured] - pieceScore[moves[i].PieceMoved]
		}
		if moves[i].IsPawnPromotion {
			score += 90
		}
		moves[i].Score = score
	}
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].Score > moves[j].Score
	})
}

func quiescenceSearch(gs *GameState, alpha, beta int, turnMultiplier int) int {
	standPat := turnMultiplier * ScoreBoard(gs)
	if standPat >= beta {
		return beta
	}
	if alpha < standPat {
		alpha = standPat
	}

	moves := gs.GetValidMoves()
	orderMoves(gs, moves)

	for _, move := range moves {
		if move.PieceCaptured != '-' {
			gs.MakeMove(move)
			score := -quiescenceSearch(gs, -beta, -alpha, -turnMultiplier)
			gs.UndoMove()

			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
	}
	return alpha
}

// Simple integer math helpers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
