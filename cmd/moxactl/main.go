package main

import (
	"fmt"
	"log"
	"os"

	"github.com/LilithGames/moxa/cli"
)

func main() {
	if err := cli.App.Run(os.Args); err != nil {
		log.Fatalln("[FATAL]", fmt.Errorf("app.Run(%v) err: %w", os.Args, err))
	}
}
