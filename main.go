package main

import (
	"fmt"

	"github.com/notnil/chess"

	"github.com/LouieMartin/polaris/engine"
)

func main() {
	var game *chess.Game = chess.NewGame()

	for game.Outcome() == chess.NoOutcome {
		engine.PlayBestMove(3, game)

		fmt.Println(game.Position().Board().Draw())
	}

	fmt.Println(game.String())
}
