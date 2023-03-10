from chess import WHITE, BLACK, Board, Move
from engine import find_best_move
from random import choice
from sys import exit

PIECE_MAP = {
    # White
    'K': '♚',
    'Q': '♛',
    'R': '♜',
    'B': '♝',
    'N': '♞',
    'P': '♟',

    # Black
    'k': '♔',
    'q': '♕',
    'r': '♖',
    'b': '♗',
    'n': '♘',
    'p': '♙',
}

def get_display(board: Board) -> str:
    """
        Display the board with files and ranks.
    """

    board_string = list(str(board))

    for index, piece in enumerate(board_string):
        if piece in PIECE_MAP:
            board_string[index] = PIECE_MAP[piece]

    display = []
    for index, rank in enumerate(''.join(board_string).split('\n')):
        display.append(f'{8 - index} {rank}')

    if board.turn == BLACK:
        display.reverse()

        for index in range(len(display)):
            display[index] = display[index][::-1]
    
    display.append(
        '  a b c d e f g h'
        if board.turn == WHITE else
        'h g f e d c b a'
    )

    return '\n' + '\n'.join(display) + '\n'

def get_move(board: Board) -> Move:
    """
        Keep getting the move from the user.
    """

    moves = list(board.legal_moves)
    move = input(f'Play a move (e.g. {choice(moves)}): ')

    for legal_move in moves:
        if str(legal_move) == move:
            return legal_move
    
    return get_move(board)

def start():
    """
        Run this function to start the command-line interface.
    """

    side = input('Play as [b]lack or [w]hite: ')
    board = Board()

    if side == 'w':
        print(get_display(board))
        move = get_move(board)

        board.push(move)
    
    while not board.is_game_over():
        move, _ = find_best_move(4, board)

        board.push(move)

        print(get_display(board))

        move = get_move(board)

        board.push(move)

try:
    start()
except KeyboardInterrupt:
    exit()
