package main

import (
	"fmt"

	"github.com/LouieMartin/polaris/engine"
	"github.com/notnil/chess"
)

func main() {
	game := chess.NewGame()
	for game.Outcome() == chess.NoOutcome {
		engine.PlayBestMove(3, game)
		fmt.Println(game.Position().Board().Draw())
	}

	fmt.Println(game.String())
}
