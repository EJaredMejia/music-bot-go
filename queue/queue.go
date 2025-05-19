package queue

import (
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

func NewQueue() *Queue {
	return &Queue{
		elements: make([]goYtldp.ProgressUpdate, 0),
	}
}

func (q *Queue) Enqueue(value goYtldp.ProgressUpdate, vc *discordgo.VoiceConnection) {
	q.mu.Lock()

	q.elements = append(q.elements, value)
	q.mu.Unlock()

	go streamAudio(q, value, vc)
}

func streamAudio(q *Queue, value goYtldp.ProgressUpdate, vc *discordgo.VoiceConnection) {
	q.audioStreamMu.Lock()
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
