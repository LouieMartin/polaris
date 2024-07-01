package engine

import (
	"math"
	"sort"

	"github.com/LouieMartin/polaris/evaluation"
	"github.com/notnil/chess"
)

func OrderMoves(moves []*chess.Move, position *chess.Position) {
	endgame := evaluation.IsEndgame(position)
	sort.Slice(moves, func(i, j int) bool {
		return evaluation.EvaluateMove(moves[i], endgame, position) < evaluation.EvaluateMove(moves[j], endgame, position)
	})

	if position.Turn() == chess.White {
		sort.SliceStable(moves, func(_, _ int) bool {
			return true
		})
	}
}

func Quiesce(alpha float64, beta float64, position *chess.Position) float64 {
	evaluation := evaluation.Evaluate(position)
	if evaluation >= beta {
		return beta
	}

	alpha = math.Max(evaluation, alpha)
	moves := position.ValidMoves()
	var captures []*chess.Move

	for _, move := range moves {
		if move.HasTag(chess.Capture) {
			captures = append(captures, move)
		}
	}

	OrderMoves(captures, position)
	for _, capture := range captures {
		newPosition := position.Update(capture)

		evaluation = -Quiesce(-beta, -alpha, newPosition)

		if evaluation >= beta {
			return beta
		}

		alpha = math.Max(evaluation, alpha)
	}

	return alpha
}

func Negamax(depth uint8, alpha float64, beta float64, position *chess.Position) float64 {
	switch position.Status() {
	case chess.Checkmate:
		return math.Inf(-1)
	case chess.Stalemate:
		return 0.0
	}

	if depth == 0 {
		return Quiesce(alpha, beta, position)
	}

	moves := position.ValidMoves()

	OrderMoves(moves, position)
	for _, move := range moves {
		newPosition := position.Update(move)
		evaluation := -Negamax(depth-1, -beta, -alpha, newPosition)

		if evaluation >= beta {
			return beta
		}

		alpha = math.Max(evaluation, alpha)
	}

	return alpha
}

func FindBestMove(depth uint8, position *chess.Position) (*chess.Move, float64) {
	switch position.Status() {
	case chess.Checkmate:
		return nil, math.Inf(-1)
	case chess.Stalemate:
		return nil, 0.0
	}

	moves := position.ValidMoves()
	bestEvaluation := math.Inf(-1)
	var bestMove *chess.Move

	OrderMoves(moves, position)

	for _, move := range moves {
		newPosition := position.Update(move)
		evaluation := -Negamax(depth-1, math.Inf(-1), math.Inf(1), newPosition)

		if bestEvaluation <= evaluation {
			bestEvaluation = evaluation
			bestMove = move
		}
	}

	return bestMove, bestEvaluation
}

func PlayBestMove(depth uint8, game *chess.Game) (*chess.Move, float64) {
	move, evaluation := FindBestMove(depth, game.Position())

	if move != nil {
		game.Move(move)
	}

	return move, evaluation
}
