package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/music-formatter/bot"
	"github.com/music-formatter/queue"
)

func main() {
	fmt.Println("main")

	// TODO use this logic for max files
	inputString := "!p --max=25"
	parts := strings.Split(inputString, " ")

	if len(parts) == 0 {
		fmt.Println("No input provided.")
		return
	}

	// action := parts[0]
	arguments := parts[1:]
	fmt.Println(arguments)

	fs := flag.NewFlagSet("command", flag.ContinueOnError)

	// var songName string
	// fs.StringVar(&songName, "song-name", "", "name of the song")
	maxSongs := fs.Int("max", 0, "max number of songs")

	err := fs.Parse(arguments)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fs.Args())

	// println(songName)
	fmt.Println(*maxSongs)

	queues := queue.CreateDynamicQueues()

	bot.Run(queues)

}
