package queue

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	goYtldp "github.com/lrstanley/go-ytdlp"

	"github.com/bwmarrin/dgvoice"
)

type DynamicQueues map[string]*Queue

func CreateDynamicQueues() DynamicQueues {
	return make(DynamicQueues)
}

func GetQueue(queues DynamicQueues, id string) *Queue {
	queue, exists := queues[id]

	if exists {
		return queue
	}

	queues[id] = NewQueue()

	return queues[id]
}

// Queue is a struct that represents a queue data structure.
// It uses a slice to store the elements and a mutex for concurrent access.
type Queue struct {
	elements      []goYtldp.ProgressUpdate
	mu            sync.Mutex
	audioStreamMu sync.Mutex
}

type Queues Queue

func NewQueue() *Queue {
	return &Queue{
		elements: make([]goYtldp.ProgressUpdate, 0),
	}
}

func (q *Queue) Enqueue(discord *discordgo.Session, discordMessage *discordgo.MessageCreate, value goYtldp.ProgressUpdate, vc *discordgo.VoiceConnection) {
	q.mu.Lock()

	q.elements = append(q.elements, value)
	q.mu.Unlock()

	go streamAudio(q, discord, discordMessage, value, vc)
}

func streamAudio(q *Queue, discord *discordgo.Session, discordMessage *discordgo.MessageCreate, value goYtldp.ProgressUpdate, vc *discordgo.VoiceConnection) {
	q.audioStreamMu.Lock()
	discord.ChannelMessageSend(discordMessage.ChannelID, fmt.Sprintf("now playing: %s", *value.Info.Title))

	dgvoice.PlayAudioFile(vc, value.Filename, make(<-chan bool))
	q.audioStreamMu.Unlock()

	q.Dequeue()
}

func getTempFilename(filename string) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + ".temp" + ext
}

func (q *Queue) Dequeue() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.IsEmpty() {
		return
	}

	element := q.elements[0]
	q.elements = q.elements[1:]
	go func() {
		tempFilename := getTempFilename(element.Filename)
		os.Remove(element.Filename)
		os.Remove(tempFilename)
	}()

	return
}

func (q *Queue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *Queue) Size() int {
	return len(q.elements)
}

func (q *Queue) Print() string {
	size := q.Size()
	if size == 0 {
		return "There is no song playing"
	}

	var sb strings.Builder
	sb.WriteString("Queue:\n")
	for i, element := range q.elements {
		isPlaying := i == 0

		message := ""

		if isPlaying {
			message = "**Currently Playing:** "
		}

		duration := formatSecondsToMMSS(*element.Info.Duration)

		sb.WriteString(fmt.Sprintf("%d. %s%s - `%s`\n", i+1, message, *element.Info.Title, duration))
	}

	return sb.String()
}

func formatSecondsToMMSS(totalSeconds float64) string {
	// Ensure totalSeconds is non-negative
	if totalSeconds < 0 {
		totalSeconds = 0
	}

	// Calculate minutes and seconds
	minutes := int(totalSeconds / 60)
	seconds := int(math.Mod(totalSeconds, 60)) // Use math.Mod for float64

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
