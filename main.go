package main

import (
	"fmt"

	"github.com/music-formatter/bot"
	"github.com/music-formatter/queue"
)

const PREFIX = "!"

func main() {
	fmt.Println("main")

	queues := queue.CreateDynamicQueues()

	bot.Run(queues)

}
