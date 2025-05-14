package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/music-formatter/queue"
	"github.com/music-formatter/ytldp"
)

const PREFIX = "!"

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
			discord: s,
			message: m,
			queues:  queues,
		})
	})

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

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

	discord.ChannelTyping(discordMessage.ChannelID)
	action, message := splitActionMessage(discordMessage.Content)
	switch {
	case isCommand(action, "br"):

		vc, err := joinVoiceChannel(discord, discordMessage)

		if err != nil {
			return
		}

		queue := queue.GetQueue(queues, vc.ChannelID)
		ytldp.ExtractAudio(queue, vc, message)

	}

}

func joinVoiceChannel(discord *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {
	voiceState, err := discord.State.VoiceState(m.GuildID, m.Author.ID)

	if err != nil {
		discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not on a voice channel", m.Author.Username))
		return nil, err
	}

	vc, err := discord.ChannelVoiceJoin(m.GuildID, voiceState.ChannelID, false, true)

	if err != nil {
		fmt.Println("Error joining voice channel:", err)
		discord.ChannelMessageSend(m.ChannelID, "Failed to join the voice channel.")
		return vc, err
	}

	return vc, err
}

func isCommand(message string, action string) bool {
	return strings.Contains(message, PREFIX+action)
}

func splitActionMessage(content string) (string, string) {
	words := strings.Split(content, " ")

	action := words[0]
	message := strings.Join(words[1:], " ")

	return action, message
}
