package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LouieMartin/polaris/engine"
	"github.com/notnil/chess"
)

func main() {
	var game *chess.Game = chess.NewGame()

	for game.Outcome() == chess.NoOutcome {
		fmt.Print("Enter your move: ")
		reader := bufio.NewReader(os.Stdin)
		move, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		if err := game.MoveStr(strings.TrimSpace(move)); err != nil {
			log.Fatal(err)
		}

		engine.PlayBestMove(3, game)

		fmt.Println(game.Position().Board().Draw())
	}

	fmt.Println(game.String())
}
