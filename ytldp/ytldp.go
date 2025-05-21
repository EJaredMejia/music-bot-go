package ytldp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	goYtldp "github.com/lrstanley/go-ytdlp"
	"github.com/music-formatter/queue"
)

type ExtractAudioParams struct {
	Queue          *queue.Queue
	Discord        *discordgo.Session
	DiscordMessage *discordgo.MessageCreate
	DiscordVc      *discordgo.VoiceConnection
	TextMessage    string
	MaxSongs       int
}

func ExtractAudio(params ExtractAudioParams) {

	discord := params.Discord
	discordMessage := params.DiscordMessage
	vc := params.DiscordVc
	queue := params.Queue
	song := params.TextMessage

	goYtldp.MustInstall(context.TODO(), nil)

	// todo see playlist end why crashes
	audioDirectory := fmt.Sprintf("audio/%s", vc.ChannelID)
	log.Println("Extracting audio")
	dl := goYtldp.
		New().
		FormatSort("res,aext").
		ExtractAudio().
		MaxDownloads(params.MaxSongs).
		ProgressFunc(100*time.Millisecond, func(progress goYtldp.ProgressUpdate) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("ERROR defer ytld: ", r)
					discord.ChannelMessageSend(discordMessage.ChannelID, fmt.Sprintf("Error: %w", r))
				}
			}()

			if !progress.Status.IsCompletedType() {
				return
			}

			log.Println("DONE")
			queue.Enqueue(discord, discordMessage, progress, vc)

		}).
		DefaultSearch("ytsearch").
		Output(audioDirectory + "/%(playlist_index)s - %(extractor)s - %(title)s.%(ext)s")

	// I could add a file so it is detected so i can remove the watcher??
	// https://www.youtube.com/watch?v=ftaXMKV3ffE

	// url := "https://www.youtube.com/watch?v=ftaXMKV3ffE"
	// url := "olympian playboi carti"

	// url := "https://www.youtube.com/watch?v=fI-mnYR-Mp8&list=RDgsbZ3KX2CR8&index=27"

	res, err := dl.Run(context.TODO(), song)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("RES: ", res)

}
