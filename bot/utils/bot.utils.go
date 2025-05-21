package utils

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const PREFIX = "!"

func IsCommand(discord *discordgo.Session, m *discordgo.MessageCreate, action string) string {
	hasPrefix := strings.Contains(action, PREFIX)

	if !hasPrefix {
		return ""
	}

	command := action[1:]

	return command
}

func JoinVoiceChannel(discord *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {
	discord.ChannelTyping(m.ChannelID)

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

func SplitActionMessage(input string) (string, string, []string) {
	var flags []string
	var messageBuilder strings.Builder
	trimInput := strings.TrimSpace(input)
	words := strings.Fields(trimInput)

	action := strings.ToLower(words[0])

	messageWords := words[1:]
	for _, word := range messageWords {

		if strings.HasPrefix(word, "--") || strings.HasPrefix(word, "-") {
			flags = append(flags, word)
			continue
		}

		messageBuilder.WriteString(word + " ")
	}

	fmt.Println("MESSAGE:", messageBuilder.String())

	return action, messageBuilder.String(), flags
}
