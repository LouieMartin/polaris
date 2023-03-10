from engine import find_best_move
from chess import Board
from sys import exit

board = Board()

while True:
    message = input().strip()
    tokens = message.split(' ')

    if message == 'isready':
        print('readyok')
    if message in ('quit', 'stop'):
        exit()

    if message == 'uci':
        print('id name Polaris')
        print('id author Louie Decierdo')
        print('uciok')

    if message[0:2] == 'go':
        move, _ = find_best_move(4, board)

        print(f'bestmove {move}')
        
    if message == 'd':
        print(board)
        print(board.fen())

    if message.startswith('position'):
        if len(tokens) < 2:
            continue
        
        if tokens[1] == 'startpos':
            moves_start = 2
            board.reset()

        if tokens[1] == 'fen':
            fen = ' '.join(tokens[2:8])

            board.set_fen(fen)
            moves_start = 8
        
        if len(tokens) <= moves_start or tokens[moves_start] != 'moves':
            continue

        for move in tokens[(moves_start + 1):]:
            board.push_uci(move)
