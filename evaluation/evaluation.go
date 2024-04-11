package evaluation

import (
	"math"

	"github.com/LouieMartin/polaris/constants"
	"github.com/notnil/chess"
)

func IsEndgame(position *chess.Position) bool {
	var queens uint8 = 0
	var minors uint8 = 0
	var square uint8 = 0
	for ; square < 64; square++ {
		var piece chess.Piece = position.Board().Piece(chess.Square(square))
		if piece.Type() == chess.Bishop || piece.Type() == chess.Knight {
			minors++
		}

		if piece.Type() == chess.Queen {
			queens++
		}
	}

	return (queens >= 2 && minors <= 1) || queens == 0
}

func GetPositionalValue(piece chess.Piece, square chess.Square, endgame bool) float64 {
	if piece.Color() == chess.Black {
		square ^= 56
	}

	if piece.Type() == chess.King && endgame {
		return constants.KingEndgameSquareTable[square]
	}

	return constants.PieceSquareTables[piece.Type()][square]
}

func EvaluatePiece(piece chess.Piece, square chess.Square, endgame bool) float64 {
	positionalValue := GetPositionalValue(piece, square, endgame)
	materialValue := constants.PieceValues[piece.Type()]
	return positionalValue + materialValue
}

func Evaluate(position *chess.Position) float64 {
	score := 0.0
	var turn float64
	endgame := IsEndgame(position)
	if turn = 1; position.Turn() == chess.Black {
		turn = -1
	}

	for square := 0; square < 64; square++ {
		piece := position.Board().Piece(chess.Square(square))
		if piece.Color() == chess.NoColor {
			continue
		}

		var color float64
		evaluation := EvaluatePiece(piece, chess.Square(square), endgame)
		if color = 1.0; piece.Color() == chess.Black {
			color = -1.0
		}

		score += evaluation * color
	}

	return turn * score
}

func EvaluateCapture(move *chess.Move, position *chess.Position) float64 {
	if move.HasTag(chess.EnPassant) {
		return constants.PieceValues[chess.Pawn]
	}

	piece := position.Board().Piece(move.S1())
	capturedPiece := position.Board().Piece(move.S2())
	return constants.PieceValues[capturedPiece.Type()] - constants.PieceValues[piece.Type()]
}

func EvaluateMove(move *chess.Move, endgame bool, position *chess.Position) float64 {
	var turn float64
	var materialChange float64
	if turn = 1; position.Turn() == chess.Black {
		turn = -1
	}

	if move.Promo() != chess.NoPieceType {
		return turn * math.Inf(1)
	}

	piece := position.Board().Piece(move.S1())
	fromValue := GetPositionalValue(piece, move.S1(), endgame)
	toValue := GetPositionalValue(piece, move.S2(), endgame)
	positionalValue := toValue - fromValue

	if move.HasTag(chess.Capture) {
		materialChange = EvaluateCapture(move, position)
	}

	return (positionalValue + materialChange) * turn
}
