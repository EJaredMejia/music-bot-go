package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/music-formatter/bot/commands"
	"github.com/music-formatter/bot/utils"
	"github.com/music-formatter/queue"
)

func Run(queues queue.DynamicQueues) {

	// create a session
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		panic(err)
	}

	// add a event handler
	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		newMessage(newMessageParams{
			queues:  queues,
			message: m,
			discord: s,
		})
	})

	// open session
	discord.Open()
	defer func() {
		os.RemoveAll("audio")
		discord.Close()
	}() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

type newMessageParams struct {
	discord *discordgo.Session
	message *discordgo.MessageCreate
	queues  queue.DynamicQueues
}

func newMessage(params newMessageParams) {
	discordMessage := params.message
	discord := params.discord
	queues := params.queues

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if discordMessage.Author.ID == discord.State.User.ID {
		return
	}

	action, message, flags := utils.SplitActionMessage(discordMessage.Content)
	commandParams := commands.CommandParams{
		Queues:         queues,
		TextMessage:    message,
		Discord:        discord,
		DiscordMessage: discordMessage,
		Flags:          flags,
	}
	switch utils.IsCommand(discord, discordMessage, action) {
	case commands.BR:
		commands.PlayCommand(commandParams)
	case commands.QUEUE:
		commands.PrintQueue(commandParams)
	case commands.LEAVE:
		commands.LeaveCommand(commandParams)
	default:
		// TODO add invalid command
	}
}
