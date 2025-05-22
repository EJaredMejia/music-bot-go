package commands

import (
	"github.com/music-formatter/bot/utils"
	"github.com/music-formatter/queue"
)

func PrintQueue(params CommandParams) {
	vc, err := utils.JoinVoiceChannel(params.Discord, params.DiscordMessage)

	if err != nil {
		return
	}

	queue := queue.GetQueue(params.Queues, vc.ChannelID)

	if queue.IsEmpty() {
		params.Discord.ChannelMessageSend(params.DiscordMessage.ChannelID, "Queue is empty")
		return
	}

	params.Discord.ChannelMessageSend(params.DiscordMessage.ChannelID, queue.Print())
}
