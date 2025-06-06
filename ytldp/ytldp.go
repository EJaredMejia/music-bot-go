package ytldp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	goYtldp "github.com/lrstanley/go-ytdlp"
	"github.com/music-formatter/queue"
)

const DEFAULT_URL_MAX_SONGS = 25
const DEFAULT_TEXT_MAX_SONGS = 1

var CANCELLED_ERROR = errors.New("Cancelled by context")

type PlayFlags struct {
	MaxSongs int
	Random   bool
}

type ExtractAudioParams struct {
	Queue          *queue.Queue
	Discord        *discordgo.Session
	DiscordMessage *discordgo.MessageCreate
	DiscordVc      *discordgo.VoiceConnection
	TextMessage    string
	Flags          PlayFlags
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func ExtractAudio(params ExtractAudioParams) {

	discord := params.Discord
	discordMessage := params.DiscordMessage
	vc := params.DiscordVc
	queue := params.Queue
	song := params.TextMessage

	goYtldp.MustInstall(context.TODO(), &goYtldp.InstallOptions{
		AllowVersionMismatch: true,
	})

	isSongUrl := isURL(song)

	audioDirectory := fmt.Sprintf("audio/%s", vc.ChannelID)
	// playlist skip -> song why cancel
	// but leave should cancel everything

	log.Println("Extracting audio")
	dl := goYtldp.
		New().
		Format("bestaudio[ext=webm]/bestaudio[ext=m4a]/bestaudio").
		FormatSort("abr").
		Cookies("cookies.txt").
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

			log.Println("DONE:", *progress.Info.Title)
			queue.Enqueue(discord, discordMessage, progress, vc)

		}).
		DefaultSearch("ytsearch").
		Output(audioDirectory + "/%(playlist_index)s - %(extractor)s - %(title)s.%(ext)s")

	if params.Flags.Random {
		dl = dl.PlaylistRandom()
	}

	maxSongs := func() int {

		maxSongs := params.Flags.MaxSongs
		isDefaultMaxSongs := maxSongs == 0

		if !isDefaultMaxSongs {
			return maxSongs
		}

		if isSongUrl {
			return DEFAULT_URL_MAX_SONGS
		}

		return DEFAULT_TEXT_MAX_SONGS
	}()

	dl = dl.MaxDownloads(maxSongs)

	done := make(chan bool, 1)

	go func() {
		defer close(done)
		res, err := dl.Run(queue.Context, song)
		if err != nil {
			log.Println("ERR:", err)
		}

		log.Println("RES: ", res)
	}()

	func() {
		select {
		case <-queue.Context.Done():
			log.Println("CANCELLED")
			// panic(errors.New("Cancelled by context"))
		case <-done:
			return
		}
	}()
	// done <- true

}
