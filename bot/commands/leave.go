package commands

import (
	"github.com/music-formatter/bot/utils"
	"github.com/music-formatter/queue"
)

func LeaveCommand(params CommandParams) {
	params.Discord.ChannelTyping(params.DiscordMessage.ChannelID)
	voiceState, err := utils.GetVoiceState(params.Discord, params.DiscordMessage)

	if err != nil {
		return
	}
	vc, ok := params.Discord.VoiceConnections[voiceState.GuildID]

	if !ok {
		params.Discord.ChannelMessageSend(params.DiscordMessage.ChannelID, "bot is not in a voice channel")
		return
	}

	// send a message to the channel
	params.Discord.ChannelMessageSend(params.DiscordMessage.ChannelID, "Left the voice channel and cleared the queue.")

	// get the queue
	queue := queue.GetQueue(params.Queues, vc.ChannelID)

	// clear the queue
	queue.Clear(params.Queues, vc)

	vc.Disconnect()
}
