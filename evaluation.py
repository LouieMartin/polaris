# This module implements Tomasz Michniewski's Simplified Evaluation Function:
#   https://www.chessprogramming.org/Simplified_Evaluation_Function

from chess import (
    # Squares, Boards, and Moves.
    SQUARES,
    Square,
    Board,
    Move,

    # Pieces.
    KING,
    QUEEN,
    ROOK,
    BISHOP,
    KNIGHT,
    PAWN,
    Piece,

    # Color.
    WHITE,
    BLACK,
)

KING_ENDGAME = 7
PIECE_VALUES = {
    KING: 20000,
    QUEEN: 900,
    ROOK: 500,
    BISHOP: 330,
    KNIGHT: 320,
    PAWN: 100,
}

# NOTE: The Piece-Square Tables are reveresed so the top-left square is A1.

WHITE_PIECE_SQUARE_TABLES = {
    KING_ENDGAME: [
        -50,-30,-30,-30,-30,-30,-30,-50,
        -30,-30,  0,  0,  0,  0,-30,-30,
        -30,-10, 20, 30, 30, 20,-10,-30,
        -30,-10, 30, 40, 40, 30,-10,-30,
        -30,-10, 30, 40, 40, 30,-10,-30,
        -30,-10, 20, 30, 30, 20,-10,-30,
        -30,-20,-10,  0,  0,-10,-20,-30,
        -50,-40,-30,-20,-20,-30,-40,-50,
    ],
    KING: [
        20, 30, 10,  0,  0, 10, 30, 20,
        20, 20,  0,  0,  0,  0, 20, 20,
        -10,-20,-20,-20,-20,-20,-20,-10,
        -20,-30,-30,-40,-40,-30,-30,-20,
        -30,-40,-40,-50,-50,-40,-40,-30,
        -30,-40,-40,-50,-50,-40,-40,-30,
        -30,-40,-40,-50,-50,-40,-40,-30,
        -30,-40,-40,-50,-50,-40,-40,-30,
    ],
    QUEEN: [
        -20,-10,-10, -5, -5,-10,-10,-20,
        -10,  0,  5,  0,  0,  0,  0,-10,
        -10,  5,  5,  5,  5,  5,  0,-10,
        0,  0,  5,  5,  5,  5,  0, -5,
        -5,  0,  5,  5,  5,  5,  0, -5,
        -10,  0,  5,  5,  5,  5,  0,-10,
        -10,  0,  0,  0,  0,  0,  0,-10,
        -20,-10,-10, -5, -5,-10,-10,-20,
    ],
    ROOK: [
        0,  0,  0,  5,  5,  0,  0,  0,
        -5,  0,  0,  0,  0,  0,  0, -5,
        -5,  0,  0,  0,  0,  0,  0, -5,
        -5,  0,  0,  0,  0,  0,  0, -5,
        -5,  0,  0,  0,  0,  0,  0, -5,
        -5,  0,  0,  0,  0,  0,  0, -5,
        5, 10, 10, 10, 10, 10, 10,  5,
        0,  0,  0,  0,  0,  0,  0,  0,
    ],
    BISHOP: [
        -20,-10,-10,-10,-10,-10,-10,-20,
        -10,  5,  0,  0,  0,  0,  5,-10,
        -10, 10, 10, 10, 10, 10, 10,-10,
        -10,  0, 10, 10, 10, 10,  0,-10,
        -10,  5,  5, 10, 10,  5,  5,-10,
        -10,  0,  5, 10, 10,  5,  0,-10,
        -10,  0,  0,  0,  0,  0,  0,-10,
        -20,-10,-10,-10,-10,-10,-10,-20,
    ],
    KNIGHT: [
        -50,-40,-30,-30,-30,-30,-40,-50,
        -40,-20,  0,  5,  5,  0,-20,-40,
        -30,  5, 10, 15, 15, 10,  5,-30,
        -30,  0, 15, 20, 20, 15,  0,-30,
        -30,  5, 15, 20, 20, 15,  5,-30,
        -30,  0, 10, 15, 15, 10,  0,-30,
        -40,-20,  0,  0,  0,  0,-20,-40,
        -50,-40,-30,-30,-30,-30,-40,-50,
    ],
    PAWN: [
        0,  0,  0,  0,  0,  0,  0,  0,
        5, 10, 10,-20,-20, 10, 10,  5,
        5, -5,-10,  0,  0,-10, -5,  5,
        0,  0,  0, 20, 20,  0,  0,  0,
        5,  5, 10, 25, 25, 10,  5,  5,
        10, 10, 20, 30, 30, 20, 10, 10,
        50, 50, 50, 50, 50, 50, 50, 50,
         0,  0,  0,  0,  0,  0,  0,  0,
    ],
}

BLACK_PIECE_SQUARE_TABLES = {
    KING_ENDGAME: WHITE_PIECE_SQUARE_TABLES[KING_ENDGAME][::-1],
    KING: WHITE_PIECE_SQUARE_TABLES[KING][::-1],
    QUEEN: WHITE_PIECE_SQUARE_TABLES[QUEEN][::-1],
    ROOK: WHITE_PIECE_SQUARE_TABLES[ROOK][::-1],
    BISHOP: WHITE_PIECE_SQUARE_TABLES[BISHOP][::-1],
    KNIGHT: WHITE_PIECE_SQUARE_TABLES[KNIGHT][::-1],
    PAWN: WHITE_PIECE_SQUARE_TABLES[PAWN][::-1],
}

PIECE_SQUARE_TABLES = {
    WHITE: WHITE_PIECE_SQUARE_TABLES,
    BLACK: BLACK_PIECE_SQUARE_TABLES,
}

def is_endgame(board: Board) -> bool:
    """
        This function checks if the game is in an endgame state according to `Per Michniewski`:
            * Both sides have queens, but there are one or less minor pieces (a bishop or knight).
            * Both sides have traded their queens.
    """

    queens = 0
    minors = 0

    for square in SQUARES:
        piece = board.piece_at(square)

        if piece:
            if piece.piece_type in (BISHOP, KNIGHT):
                minors += 1
            elif piece.piece_type == QUEEN:
                queens += 1
    
    return queens == 0 or (queens >= 2 and minors <= 1)

def get_positional_value(square: Square, piece: Piece, endgame: bool) -> int:
    """
        Retrieves the positional value from the Piece-Square Tables.
    """
    
    piece_square_tables = PIECE_SQUARE_TABLES[piece.color]

    if piece.piece_type == KING and endgame:
        return piece_square_tables[KING_ENDGAME][square]
    
    return piece_square_tables[piece.piece_type][square]

def evaluate_capture(move: Move, board: Board) -> int:
    """
        Evaluates a capture by subtracting the captured piece value by the piece value.
    """

    if not board.is_capture(move):
        raise Exception('Expected a capture.')
    
    if board.is_en_passant(move):
        return PIECE_VALUES[PAWN]
    
    captured_piece = board.piece_at(move.to_square)
    piece = board.piece_at(move.from_square)

    if not captured_piece and not piece:
        raise Exception(f'Expected pieces on squares: {move.to_square} and {move.from_square}.')
    elif not captured_piece:
        raise Exception(f'Expected a piece on the {move.to_square} square.')
    elif not piece:
        raise Exception(f'Expected a piece on the {move.from_square} square.')

    return PIECE_VALUES[captured_piece.piece_type] - PIECE_VALUES[piece.piece_type]

def evaluate_move(move: Move, endgame: bool, board: Board) -> float:
    """
        Evaluates a move by adding the positional change by the material change and multiplying it by the turn to move.
    """

    perspective = 1 if board.turn == WHITE else -1
    material_change = 0

    if move.promotion:
        return perspective * float('inf')

    piece = board.piece_at(move.from_square)

    if piece:
        from_value = get_positional_value(move.from_square, piece, endgame)
        to_value = get_positional_value(move.to_square, piece, endgame)

        positional_change = to_value - from_value
    else:
        raise Exception(f'Expected a piece on the {move.from_square} square.')

    if board.is_capture(move):
        material_change = evaluate_capture(move, board)
    
    return (positional_change + material_change) * perspective

def evaluate_piece(square: Square, endgame: bool, board: Board) -> int:
    """
        Evaluates a piece by adding the positional value by the material value.
    """

    piece = board.piece_at(square)
    
    if not piece:
        raise Exception(f'Expected a piece on the {square} square.')

    positional_value = get_positional_value(square, piece, endgame)
    material_value = PIECE_VALUES[piece.piece_type]

    return positional_value + material_value

def evaluate(board: Board) -> float:
    """
        Evaluates the board by evaluating every single piece on the board with the `evaluate_piece` function.
    """

    perspective = 1 if board.turn == WHITE else -1
    endgame = is_endgame(board)
    score = 0.0

    for square in SQUARES:
        piece = board.piece_at(square)

        if not piece: continue
        evaluation = evaluate_piece(square, endgame, board)
        color = 1 if piece.color == WHITE else -1
        score += evaluation * color
    
    return perspective * score
