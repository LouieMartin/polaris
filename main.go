package main

import (
	"fmt"

	"github.com/notnil/chess"
)

func main() {
	var game *chess.Game = chess.NewGame()

	for game.Outcome() == chess.NoOutcome {
		playBestMove(3, game)

		fmt.Println(game.Position().Board().Draw())
	}

	fmt.Println(game.String())
}
