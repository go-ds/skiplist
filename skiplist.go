package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	defaultProb     float64 = 0.5
	defaultMaxLevel int     = 17
	minLevel        int     = 1
	maxLevel        int     = 64
)

var maxLevelErr = fmt.Errorf("maxLevel for a SkipList must between [%d, %d]", minLevel, maxLevel)

type node struct {
	next  []*node
	prev  []*node
	key   float64
	value interface{}
}

// SkipList implements a skip list structure.
// All operations are concurrency safe.
type SkipList struct {
	head       *node
	maxLevel   int
	level      int
	prob       float64
	probs      []float64
	length     int
	randSource rand.Source
	mut        sync.RWMutex
}

// New creates a new skip list instance.
func New(opts ...Option) *SkipList {
	list := &SkipList{
		head:       &node{},
		maxLevel:   defaultMaxLevel,
		level:      1,
		prob:       defaultProb,
		randSource: rand.NewSource(time.Now().UnixNano()),
	}

	for _, opt := range opts {
		opt(list)
	}

	list.head = &node{next: make([]*node, list.maxLevel, list.maxLevel)}

	list.makeProbs()

	return list
}

// Search finds a node by key. It returns the node value if found or nil.
func (list *SkipList) Search(key float64) interface{} {
	list.mut.RLock()
	defer list.mut.RUnlock()

	if n := list.find(key); n != nil {
		return n.value
	}

	return nil
}

// Insert adds a value into the list with the specified key.
// it updates the node value if the key exists.
func (list *SkipList) Insert(key float64, value interface{}) {
	list.mut.Lock()
	defer list.mut.Unlock()

	var n *node
	// Fast path, to see if key exists.
	if n = list.find(key); n != nil {
		// Find node, then update the value.
		n.value = value
		return
	}

	// Get level for new node.
	level := list.randLevel()
	n = &node{
		next:  make([]*node, level, level),
		prev:  make([]*node, level, level),
		key:   key,
		value: value,
	}

	cur := list.head
	// Update every level list
	for i := level - 1; i >= 0; i-- {
		for cur.next[i] != nil && cur.next[i].key < key {
			cur = cur.next[i]
		}
		if cur.next[i] != nil {
			cur.next[i].prev[i] = n
		}
		n.prev[i] = cur
		n.next[i] = cur.next[i]
		cur.next[i] = n
	}

	if level > list.level {
		// Update list level.
		list.level = level
	}

	list.length++
}

func (list *SkipList) find(key float64) *node {
	cur := list.head
	for i := list.level - 1; i >= 0; i-- {
		for cur.next[i] != nil && cur.next[i].key < key {
			cur = cur.next[i]
		}

		if cur.next[i] != nil && cur.next[i].key == key {
			return cur.next[i]
		}
	}
	return nil
}

// Delete removes a node by key from the list.
func (list *SkipList) Delete(key float64) {
	_ = list.Pop(key)
}

// Pop removes a node by key from the list.
// It returns that node value if found or nil.
func (list *SkipList) Pop(key float64) interface{} {
	list.mut.Lock()
	defer list.mut.Unlock()

	var n *node
	// Fast path, to see if key exists.
	if n = list.find(key); n == nil {
		return nil
	}

	level := len(n.next)
	for i := level - 1; i >= 0; i-- {
		n.prev[i].next[i] = n.next[i]
		if n.next[i] != nil {
			n.next[i].prev[i] = n.prev[i]
		}
	}

	if level == list.level {
		// Try to decrease level.
		cur := list.head
		for i := level - 1; i >= 1; i-- {
			// No more nodes in this level.
			if cur.next[i] == nil {
				list.level--
			}
		}
	}

	list.length--

	return n.value
}

// Len returns length of the skip list.
func (list *SkipList) Len() int {
	list.mut.Lock()
	defer list.mut.Unlock()

	return list.length
}

// String returns list info
func (list *SkipList) String() string {
	var sb strings.Builder

	for i := 0; i < list.level; i++ {
		cur := list.head
		sb.WriteString(fmt.Sprintf("level %2d", i+1))
		prev := false
		for cur.next[i] != nil {
			cur = cur.next[i]
			if prev {
				sb.WriteString(" <")
			} else {
				sb.WriteString(" ")
				prev = true
			}
			sb.WriteString("--> ")
			sb.WriteString(fmt.Sprintf("%f(%v)", cur.key, cur.value))
		}
		sb.WriteString(" --> nil\n")
	}

	return sb.String()
}

const maxRand float64 = 1 << 63

func (list *SkipList) randLevel() (lvl int) {
	r := float64(list.randSource.Int63()) / maxRand
	for lvl = 1; lvl < list.maxLevel && r < list.probs[lvl]; lvl++ {
	}
	return
}

func (list *SkipList) makeProbs() {
	list.probs = make([]float64, list.maxLevel, list.maxLevel)
	for i := 1; i < list.maxLevel; i++ {
		list.probs[i] = math.Pow(list.prob, float64(i))
	}
}

// Option specifies an option for skip list.
type Option func(list *SkipList)

// WithMaxLevel specifies the max level for skip list.
// It panics if max level isn't between [1, 64].
func WithMaxLevel(level int) Option {
	return func(list *SkipList) {
		if level < minLevel || level > maxLevel {
			panic(maxLevelErr)
		}
		list.maxLevel = level
	}
}

// WithProb specifies the probability for skip list.
func WithProb(prob float64) Option {
	return func(list *SkipList) {
		list.prob = prob
	}
}

// WithRandSource specifies the rand source for skip list.
func WithRandSource(randSource rand.Source) Option {
	return func(list *SkipList) {
		list.randSource = randSource
	}
}
