package main

import (
	"math"

	"github.com/notnil/chess"
)

func isEndgame(position *chess.Position) bool {
	var queens uint8 = 0
	var minors uint8 = 0

	for square := 0; square < 64; square++ {
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

func getPositionalValue(piece chess.Piece, square chess.Square, endgame bool) float64 {
	var pieceSquareTables [6][64]float64 = pieceSquareTables[piece.Color()-1]

	if piece.Type() == chess.King && endgame {
		if piece.Color() == chess.White {
			return whiteKingEndgameSquareTable[square]
		} else {
			return blackKingEndgameSquareTable[square]
		}
	}

	return pieceSquareTables[piece.Type()-1][square]
}

func evaluatePiece(piece chess.Piece, square chess.Square, endgame bool) float64 {
	var positionalValue float64 = getPositionalValue(piece, square, endgame)
	var materialValue float64 = pieceValues[piece.Type()-1]

	return positionalValue + materialValue
}

func evaluate(position *chess.Position) float64 {
	var endgame bool = isEndgame(position)
	var perspective float64
	var score float64 = 0.0

	if perspective = 1; position.Turn() == chess.Black {
		perspective = -1
	}

	for square := 0; square < 64; square++ {
		var piece chess.Piece = position.Board().Piece(chess.Square(square))

		if piece.Color() == chess.NoColor {
			continue
		}

		var evaluation = evaluatePiece(piece, chess.Square(square), endgame)
		var color float64

		if color = 1; piece.Color() == chess.Black {
			color = -1
		}

		score += evaluation * color
	}

	return perspective * score
}

func evaluateCapture(move *chess.Move, position *chess.Position) float64 {
	if move.HasTag(chess.EnPassant) {
		return pieceValues[chess.Pawn-1]
	}

	var capturedPiece chess.Piece = position.Board().Piece(move.S2())
	var piece chess.Piece = position.Board().Piece(move.S1())

	return pieceValues[capturedPiece.Type()-1] - pieceValues[piece.Type()-1]
}

func evaluateMove(move *chess.Move, endgame bool, position *chess.Position) float64 {
	var materialChange float64
	var perspective float64

	if perspective = 1; position.Turn() == chess.Black {
		perspective = -1
	}

	if move.Promo() != chess.NoPieceType {
		return perspective * math.Inf(1)
	}

	var piece chess.Piece = position.Board().Piece(move.S1())
	var fromValue float64 = getPositionalValue(piece, move.S1(), endgame)
	var toValue = getPositionalValue(piece, move.S2(), endgame)
	var positionalValue = toValue - fromValue

	if move.HasTag(chess.Capture) {
		materialChange = evaluateCapture(move, position)
	}

	return (positionalValue + materialChange) * perspective
}
