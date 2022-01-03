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
func (q *Queue) Add(songs []*music.Song, userId string, shuffle bool, policy AddPolicy) []*QueueItem {
	q.mu.Lock()
	items := make([]*QueueItem, len(songs))
	for i, song := range songs {
		items[i] = &QueueItem{song, userId}
	}
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	}
	list := make([]*QueueItem, len(q.List)+len(items))
	switch policy {
	case Next, Now:
		for i := 0; i <= q.Index; i++ {
			list[i] = q.List[i]
		}
		for i := 0; i < len(items); i++ {
			list[q.Index+1+i] = items[i]
		}
		for i := q.Index + 1; i < len(q.List); i++ {
			list[i+len(items)] = q.List[i]
		}
		q.Dynamic = 0
	case Last:
		for i := 0; i < len(q.List); i++ {
			list[i] = q.List[i]
		}
		for i := 0; i < len(items); i++ {
			list[len(q.List)+i] = items[i]
		}
		q.Dynamic = 0
	case Smart:
		start := max(q.Index, 0)
		for i := start; i < q.Dynamic; i++ {
			if q.List[i].QueuedBy == userId {
				start = i
			}
		}
		for i := 0; i < start; i++ {
			list[i] = q.List[i]
		}
		count := map[string]int{}
		for j, item := range items {
			for start < q.Dynamic && count[q.List[start].QueuedBy] <= j {
				list[start+j] = q.List[start]
				count[q.List[start].QueuedBy]++
				start++
			}
			list[start+j] = item
		}
		for i := start; i < len(q.List); i++ {
			list[i+len(items)] = q.List[i]
		}
		q.Dynamic = max(q.Index+1, q.Dynamic) + len(items)
	}
	q.List = list
	if policy == Now || q.Index == -1 {
		q.Index = q.Index + 1
		q.update()
	}
	q.mu.Unlock()
	return items
}

func (q *Queue) SkipTo(index int) (*QueueItem, error) {
	q.mu.Lock()
	if index >= len(q.List) || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	q.Dynamic = 0
	q.Index = index
	q.update()
	q.mu.Unlock()
	return q.List[index], nil
}

func (q *Queue) Move(from, to int) (*QueueItem, error) {
	q.mu.Lock()
	if from >= len(q.List) || to >= len(q.List) || from < 0 || to < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	q.Dynamic = 0
	target := q.List[from]
	q.List = append(q.List[:from], q.List[from+1:]...)
	q.mu.Unlock()
	return target, nil
}

func (q *Queue) Remove(index int) (*QueueItem, error) {
	q.mu.Lock()
	if index >= len(q.List) || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	target := q.List[index]
	q.List = append(q.List[:index], q.List[index+1:]...)
	if index == q.Index {
		q.update()
	} else if index < q.Index {
		q.Index -= 1
	}
	q.mu.Unlock()
	return target, nil
}

// Clears the queue
func (q *Queue) Clear() {
	q.List = []*QueueItem{}
	q.Index = -1
	q.update()
}

// Shuffles all items from the current index+1 to the end of the queue.
func (q *Queue) Shuffle() {
	q.mu.Lock()
	if len(q.List) == 0 {
		return
	}
	offset := q.Index + 1
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q.List)-offset, func(i, j int) {
		q.List[i+offset], q.List[j+offset] = q.List[j+offset], q.List[i+offset]
	})
	q.mu.Unlock()
}

// Loads the next item in the queue determined by policy defined by the
// NextPolicy member of Queue.
func (q *Queue) NextSong() {
	q.mu.Lock()
	if len(q.List) == 0 {
		return
	}
	switch q.NextPolicy {
	case LoopSong:
		break
	case LoopQueue:
		q.Index = (q.Index + 1) % len(q.List)
	case NoLoop:
		q.Index = q.Index + 1
		if q.Index >= len(q.List) {
			q.Index = -1
		}
	}
	q.update()
	q.mu.Unlock()
}

// Returns the current item in the queue that should be being played. Returns
// nil if no such item exists.
func (q *Queue) CurrentItem() *QueueItem {
	if q.Index == -1 {
		return nil
	}
	return q.List[q.Index]
}

// This is a utility function that calls the UpdateHandler based on the value
// of the current Index in the queue
func (q *Queue) update() {
	if q.UpdateHandler == nil {
		return
	}
	if q.CurrentItem() == nil {
		q.UpdateHandler(nil)
		return
	}
	q.UpdateHandler(q.CurrentItem().Song)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
