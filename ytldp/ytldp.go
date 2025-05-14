package ytldp

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	goYtldp "github.com/lrstanley/go-ytdlp"
	"github.com/music-formatter/queue"
)

func ExtractAudio(queue *queue.Queue, vc *discordgo.VoiceConnection, song string) {
	goYtldp.MustInstall(context.TODO(), nil)

	log.Println("Extracting audio")
	dl := goYtldp.
		New().
		FormatSort("res,aext").
		ExtractAudio().
		ProgressFunc(100*time.Millisecond, func(progress goYtldp.ProgressUpdate) {
			isDone := progress.Status.IsCompletedType()

			if isDone {
				queue.Enqueue(progress, vc)
			}
		}).
		DefaultSearch("ytsearch").
		// DumpJSON()
		// TODO unique identifier for request
		Output("audio/%(playlist_index)s - %(extractor)s - %(title)s.%(ext)s")

	// I could add a file so it is detected so i can remove the watcher??
	// https://www.youtube.com/watch?v=ftaXMKV3ffE

	// url := "https://www.youtube.com/watch?v=ftaXMKV3ffE"
	// url := "olympian playboi carti"

	// url := "https://www.youtube.com/watch?v=fI-mnYR-Mp8&list=RDgsbZ3KX2CR8&index=27"

	res, err := dl.Run(context.TODO(), song)
	// info, err := res.GetExtractedInfo()
	// // log.Println("INFO: ", info)
	// j, _ := json.MarshalIndent(info, "", "    ")
	// log.Println(string(j))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RES: ", res)

}
