package commands

import (
	"flag"

	"github.com/music-formatter/bot/utils"
	"github.com/music-formatter/queue"
	"github.com/music-formatter/ytldp"
)

func PlayCommand(params CommandParams) {
	vc, err := utils.JoinVoiceChannel(params.Discord, params.DiscordMessage)

	if err != nil {
		return
	}

	queue := queue.GetQueue(params.Queues, vc.ChannelID)

	flags := playFlags(params.Flags)

	ytldp.ExtractAudio(ytldp.ExtractAudioParams{
		Queue:          queue,
		TextMessage:    params.TextMessage,
		Discord:        params.Discord,
		DiscordMessage: params.DiscordMessage,
		DiscordVc:      vc,
		Flags:          flags,
	})
}

func playFlags(flags []string) ytldp.PlayFlags {
	fs := flag.NewFlagSet("command", flag.ContinueOnError)

	maxSongs := fs.Int("max", 0, "max number of songs")
	random := fs.Bool("random", false, "play songs in random order")

	fs.Parse(flags)

	return ytldp.PlayFlags{
		MaxSongs: *maxSongs,
		Random:   *random,
	}
}
