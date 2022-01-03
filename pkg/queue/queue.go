package queue

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/yjp20/orpheus/pkg/music"
)

type Queue struct {
	Index         int
	List          []*QueueItem
	NextPolicy    NextPolicy
	Dynamic       int // Exclusive right-bound for the index of songs that are "dynamic"
	UpdateHandler func(song *music.Song)

	mu sync.Mutex
}

type QueueItem struct {
	Song     *music.Song `json:"song"`
	QueuedBy string      `json:"queued_by"`
}

type AddPolicy int

const (
	Next AddPolicy = iota
	Now
	Last
	Smart
)

type NextPolicy int

const (
	LoopQueue NextPolicy = iota
	LoopSong
	NoLoop
)

func NewQueue() Queue {
	return Queue{
		Index:         -1,
		List:          []*QueueItem{},
		NextPolicy:    NoLoop,
		Dynamic:       0,
	}
}

// Adds multiple songs to the queue based on policy.
func (queue *Queue) Add(songs []*music.Song, userId string, shuffle bool, policy AddPolicy) []*QueueItem {
	queue.mu.Lock()
	items := make([]*QueueItem, len(songs))
	for i, song := range songs {
		items[i] = &QueueItem{song, userId}
	}
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	}
	list := make([]*QueueItem, len(queue.List)+len(items))
	switch policy {
	case Next, Now:
		for i := 0; i <= queue.Index; i++ {
			list[i] = queue.List[i]
		}
		for i := 0; i < len(items); i++ {
			list[queue.Index+1+i] = items[i]
		}
		for i := queue.Index + 1; i < len(queue.List); i++ {
			list[i+len(items)] = queue.List[i]
		}
		queue.Dynamic = 0
	case Last:
		for i := 0; i < len(queue.List); i++ {
			list[i] = queue.List[i]
		}
		for i := 0; i < len(items); i++ {
			list[len(queue.List)+i] = items[i]
		}
		queue.Dynamic = 0
	case Smart:
		start := max(queue.Index, 0)
		for i := start; i < queue.Dynamic; i++ {
			if queue.List[i].QueuedBy == userId {
				start = i
			}
		}
		for i := 0; i < start; i++ {
			list[i] = queue.List[i]
		}
		count := map[string]int{}
		for j, item := range items {
			for start < queue.Dynamic && count[queue.List[start].QueuedBy] <= j {
				list[start+j] = queue.List[start]
				count[queue.List[start].QueuedBy]++
				start++
			}
			list[start+j] = item
		}
		for i := start; i < len(queue.List); i++ {
			list[i+len(items)] = queue.List[i]
		}
		queue.Dynamic = max(queue.Index+1, queue.Dynamic) + len(items)
	}
	queue.List = list
	if policy == Now || queue.Index == -1 {
		queue.Index = queue.Index + 1
		queue.update()
	}
	queue.mu.Unlock()
	return items
}

func (queue *Queue) SkipTo(index int) (*QueueItem, error) {
	queue.mu.Lock()
	if index >= len(queue.List) || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	queue.Dynamic = 0
	queue.Index = index
	queue.update()
	queue.mu.Unlock()
	return queue.List[index], nil
}

func (queue *Queue) Move(from, to int) (*QueueItem, error) {
	queue.mu.Lock()
	if from >= len(queue.List) || to >= len(queue.List) || from < 0 || to < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	queue.Dynamic = 0
	target := queue.List[from]
	queue.List = append(queue.List[:from], queue.List[from+1:]...)
	queue.mu.Unlock()
	return target, nil
}

func (queue *Queue) Remove(index int) (*QueueItem, error) {
	queue.mu.Lock()
	if index >= len(queue.List) || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	target := queue.List[index]
	queue.List = append(queue.List[:index], queue.List[index+1:]...)
	if index == queue.Index {
		queue.update()
	} else if index < queue.Index {
		queue.Index -= 1
	}
	queue.mu.Unlock()
	return target, nil
}

// Clears the queue
func (queue *Queue) Clear() {
	queue.List = []*QueueItem{}
	queue.Index = -1
	queue.update()
}

// Shuffles all items from the current index+1 to the end of the queue.
func (queue *Queue) Shuffle() {
	queue.mu.Lock()
	if len(queue.List) == 0 {
		return
	}
	offset := queue.Index + 1
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(queue.List)-offset, func(i, j int) {
		queue.List[i+offset], queue.List[j+offset] = queue.List[j+offset], queue.List[i+offset]
	})
	queue.mu.Unlock()
}

// Loads the next item in the queue determined by policy defined by the
// NextPolicy member of Queue.
func (queue *Queue) NextSong() {
	queue.mu.Lock()
	if len(queue.List) == 0 {
		return
	}
	switch queue.NextPolicy {
	case LoopSong:
		break
	case LoopQueue:
		queue.Index = (queue.Index + 1) % len(queue.List)
	case NoLoop:
		queue.Index = queue.Index + 1
		if queue.Index >= len(queue.List) {
			queue.Index = -1
		}
	}
	queue.update()
	queue.mu.Unlock()
}

// Returns the current item in the queue that should be being played. Returns
// nil if no such item exists.
func (queue *Queue) CurrentItem() *QueueItem {
	if queue.Index == -1 {
		return nil
	}
	return queue.List[queue.Index]
}

// This is a utility function that calls the UpdateHandler based on the value
// of the current Index in the queue
func (queue *Queue) update() {
	if queue.UpdateHandler == nil {
		return
	}
	if queue.CurrentItem() == nil {
		queue.UpdateHandler(nil)
		return
	}
	queue.UpdateHandler(queue.CurrentItem().Song)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
