package main

import (
	"flag"
	"sword/game"
)

func main() {
	usePlaceholders := flag.Bool("placeholders", false, "Use placeholder sprites instead of actual sprites")
	flag.Parse()
	if err := game.Run(*usePlaceholders); err != nil {
		panic(err)
	}
}