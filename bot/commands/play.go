package commands

import (
	"flag"

	"github.com/music-formatter/bot/utils"
	"github.com/music-formatter/queue"
	"github.com/music-formatter/ytldp"
)

const DEFAULT_MAX_SONGS = 25

func PlayCommand(params CommandParams) {
	vc, err := utils.JoinVoiceChannel(params.Discord, params.DiscordMessage)

	if err != nil {
		return
	}

	queue := queue.GetQueue(params.Queues, vc.ChannelID)

	maxSongs := playFlags(params.Flags)

	ytldp.ExtractAudio(ytldp.ExtractAudioParams{
		Queue:          queue,
		TextMessage:    params.TextMessage,
		Discord:        params.Discord,
		DiscordMessage: params.DiscordMessage,
		DiscordVc:      vc,
		MaxSongs:       maxSongs,
	})
}

func playFlags(flags []string) int {
	fs := flag.NewFlagSet("command", flag.ContinueOnError)

	maxSongs := fs.Int("max", DEFAULT_MAX_SONGS, "max number of songs")

	fs.Parse(flags)

	return *maxSongs
}
