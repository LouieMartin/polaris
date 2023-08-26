package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/notnil/chess"
)

func orderMoves(moves []*chess.Move, position *chess.Position) {
	var endgame = isEndgame(position)

	sort.Slice(moves, func(i, j int) bool {
		return evaluateMove(moves[i], endgame, position) < evaluateMove(moves[j], endgame, position)
	})

	if position.Turn() == chess.White {
		sort.SliceStable(moves, func(_, _ int) bool {
			return true
		})
	}
}

func quiesce(alpha float64, beta float64, position *chess.Position) float64 {
	var evaluation float64 = evaluate(position)

	if evaluation >= beta {
		return beta
	}

	alpha = math.Max(evaluation, alpha)
	var moves []*chess.Move = position.ValidMoves()
	var captures []*chess.Move

	for _, move := range moves {
		if move.HasTag(chess.Capture) {
			captures = append(captures, move)
		}
	}

	orderMoves(captures, position)

	for _, capture := range captures {
		var newPosition *chess.Position = position.Update(capture)

		evaluation = -quiesce(-beta, -alpha, newPosition)

		if evaluation >= beta {
			return beta
		}

		alpha = math.Max(evaluation, alpha)
	}

	return alpha
}

func negamax(depth uint8, alpha float64, beta float64, position *chess.Position) float64 {
	var moves []*chess.Move = position.ValidMoves()

	if position.Status() == chess.Checkmate {
		return math.Inf(-1)
	}

	if len(moves) == 0 {
		return 0.0
	}

	if depth == 0 {
		return quiesce(alpha, beta, position)
	}

	orderMoves(moves, position)

	for _, move := range moves {
		var newPosition *chess.Position = position.Update(move)
		var evaluation float64 = -negamax(depth-1, -beta, -alpha, newPosition)

		if evaluation >= beta {
			return beta
		}

		alpha = math.Max(evaluation, alpha)
	}

	return alpha
}

func findBestMove(depth uint8, game *chess.Game) (*chess.Move, float64) {
	var position = game.Position()
	var moves []*chess.Move = position.ValidMoves()

	if position.Status() == chess.Checkmate {
		return nil, math.Inf(-1)
	}

	if len(moves) == 0 {
		return nil, 0.0
	}

	var bestEvaluation float64 = math.Inf(-1)
	var bestMove *chess.Move

	orderMoves(moves, position)

	for _, move := range moves {
		var newPosition = position.Update(move)
		var evaluation float64 = -negamax(depth-1, math.Inf(-1), math.Inf(1), newPosition)

		if bestEvaluation <= evaluation {
			bestEvaluation = evaluation
			bestMove = move
		}
	}

	return bestMove, bestEvaluation
}

func playBestMove(depth uint8, game *chess.Game) (*chess.Move, float64) {
	var move, evaluation = findBestMove(depth, game)

	if move != nil {
		game.Move(move)
	}

	return move, evaluation
}
