package commands

import (
	"github.com/music-formatter/bot/utils"
	"github.com/music-formatter/queue"
)

func SkipCommand(params CommandParams) {
	vc, err := utils.JoinVoiceChannel(params.Discord, params.DiscordMessage)

	if err != nil {
		return
	}

	currentQueue := queue.GetQueue(params.Queues, vc.ChannelID)

	if currentQueue.IsEmpty() {
		params.Discord.ChannelMessageSend(params.DiscordMessage.ChannelID, "Queue is empty")
		return
	}

	song := currentQueue.Skip()

	params.Discord.ChannelMessageSend(params.DiscordMessage.ChannelID, "Skipped song: "+*song.Info.Title)
}
