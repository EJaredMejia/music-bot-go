package queue

import (
	"log"
	"os"
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
	elements []goYtldp.ProgressUpdate
	mu       sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		elements: make([]goYtldp.ProgressUpdate, 0),
	}
}

func (q *Queue) Enqueue(value goYtldp.ProgressUpdate, vc *discordgo.VoiceConnection) {
	q.mu.Lock()
	log.Println("ADDED QUEUE: ", value)

	q.elements = append(q.elements, value)
	q.mu.Unlock()

	go streamAudio(q, value, vc)
}

func streamAudio(q *Queue, value goYtldp.ProgressUpdate, vc *discordgo.VoiceConnection) {
	// TODO on enqueue stream audio
	log.Println("PLAY AUDIO")
	dgvoice.PlayAudioFile(vc, value.Filename, make(<-chan bool))
	log.Println("END PLAY AUDIO")
	// streamaudio.PlayAudio(vc, value.Filename)
	q.Dequeue()
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
		os.Remove(element.Filename)
	}()

	return
}

func (q *Queue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *Queue) Size() int {
	return len(q.elements)
}
