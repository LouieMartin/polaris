from evaluation import (
    evaluate_move,
    is_endgame,
    evaluate,
)

from chess import (
    # Squares, Boards, and Moves.
    Board,
    Move,

    # Color.
    WHITE
)

def get_ordered_moves(captures: bool, board: Board) -> list[Move]:
    """
        Retrieves the ordered moves with the `evaluate_move` function.
    """

    endgame = is_endgame(board)

    ordered_moves = sorted(
        board.legal_moves,
        reverse=(board.turn == WHITE),
        key=lambda move: evaluate_move(move, endgame, board),
    )

    if captures:
        return [move for move in ordered_moves if board.is_capture(move)]

    return ordered_moves

def quiesce(alpha: float, beta: float, board: Board) -> float:
    """
        This function loops over all the captures until no captures are left, then returns the evaluation of the position with no captures.
    """

    evaluation = evaluate(board)
    if evaluation >= beta:
        return beta

    alpha = max(evaluation, alpha)
    captures = get_ordered_moves(True, board)

    for capture in captures:
        board.push(capture)
        evaluation = -quiesce(-beta, -alpha, board)

        board.pop()
        if evaluation >= beta:
            return beta
        
        alpha = max(evaluation, alpha)
    
    return alpha

def negamax(
    depth: int,
    alpha: float,
    beta: float,
    board: Board,
) -> float:
    """
        Evaluates a position to a specific depth with Minimax, Alpha-Beta Pruning and Move Ordering.
    """

    depth = max(depth, 0)
    if board.is_checkmate():
        return float('-inf')
    elif board.is_game_over():
        return 0.0
    
    if depth == 0:
        return quiesce(alpha, beta, board)

    moves = get_ordered_moves(False, board)

    for move in moves:
        board.push(move)
        evaluation = -negamax(
            depth - 1,
            -beta,
            -alpha,
            board
        )

        board.pop()
        if evaluation >= beta:
            return beta
        
        alpha = max(evaluation, alpha)

    return alpha

def find_best_move(
    depth: int,
    board: Board,
) -> tuple[Move, float]:
    """
        Finds the best move with the `negamax` function.
    """

    moves = get_ordered_moves(False, board)
    best_evaluation = float('-inf')

    for move in moves:
        board.push(move)
        if board.can_claim_draw():
            evaluation = 0.0
        else:
            evaluation = -negamax(
                depth - 1,
                float('-inf'),
                float('inf'),
                board,
            )

        board.pop()
        if best_evaluation <= evaluation:
            best_evaluation = evaluation
            best_move = move

    return best_move, best_evaluation
