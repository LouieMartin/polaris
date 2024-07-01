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
<<<<<<< HEAD
	score := evaluation.Evaluate(position)
	alpha = math.Max(alpha, score)
	if alpha >= beta {
		return alpha
	}

=======
	evaluation := evaluation.Evaluate(position)

	if evaluation >= beta {
		return beta
	}

	alpha = math.Max(evaluation, alpha)
>>>>>>> parent of c6c1530 (Implemented NegaScout)
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

<<<<<<< HEAD
func Negascout(depth uint8, alpha float64, beta float64, position *chess.Position) float64 {
	switch position.Status() {
	case chess.Checkmate:
=======
func Negamax(depth uint8, alpha float64, beta float64, position *chess.Position) float64 {
	moves := position.ValidMoves()

	if position.Status() == chess.Checkmate {
>>>>>>> parent of c6c1530 (Implemented NegaScout)
		return math.Inf(-1)
	}

	if len(moves) == 0 {
		return 0.0
	}

	if depth == 0 {
		return Quiesce(alpha, beta, position)
	}

	OrderMoves(moves, position)

	for _, move := range moves {
		newPosition := position.Update(move)
<<<<<<< HEAD
		score := -Negascout(depth-1, -b, -alpha, newPosition)
		if score > alpha && score < beta && index > 0 {
			score = -Negascout(depth-1, -beta, -alpha, newPosition)
=======
		evaluation := -Negamax(depth-1, -beta, -alpha, newPosition)

		if evaluation >= beta {
			return beta
>>>>>>> parent of c6c1530 (Implemented NegaScout)
		}

		alpha = math.Max(evaluation, alpha)
	}

	return alpha
}

func FindBestMove(depth uint8, game *chess.Game) (*chess.Move, float64) {
	position := game.Position()
	moves := position.ValidMoves()

	if position.Status() == chess.Checkmate {
		return nil, math.Inf(-1)
	}

	if len(moves) == 0 {
		return nil, 0.0
	}

<<<<<<< HEAD
	var bestMove *chess.Move
	bestScore := math.Inf(-1)
	moves := position.ValidMoves()
=======
	bestEvaluation := math.Inf(-1)
	var bestMove *chess.Move

>>>>>>> parent of c6c1530 (Implemented NegaScout)
	OrderMoves(moves, position)

	for _, move := range moves {
		newPosition := position.Update(move)
<<<<<<< HEAD
		score := -Negascout(depth-1, math.Inf(-1), math.Inf(1), newPosition)
		if bestScore <= score {
			bestScore = score
=======
		evaluation := -Negamax(depth-1, math.Inf(-1), math.Inf(1), newPosition)

		if bestEvaluation <= evaluation {
			bestEvaluation = evaluation
>>>>>>> parent of c6c1530 (Implemented NegaScout)
			bestMove = move
		}
	}

	return bestMove, bestEvaluation
}

func PlayBestMove(depth uint8, game *chess.Game) (*chess.Move, float64) {
	move, evaluation := FindBestMove(depth, game)

	if move != nil {
		game.Move(move)
	}

	return move, evaluation
}
