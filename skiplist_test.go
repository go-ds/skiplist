package skiplist

import (
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMaxLevel(t *testing.T) {
	t.Parallel()

	t.Run("less than min level", func(t *testing.T) {
		assert.Panics(t, func() {
			New(WithMaxLevel(0))
		})
	})

	t.Run("big than max level", func(t *testing.T) {
		assert.Panics(t, func() {
			New(WithMaxLevel(65))
		})
	})

	t.Run("set custom max level", func(t *testing.T) {
		list := New(WithMaxLevel(11))
		assert.Equal(t, 11, list.maxLevel)
	})
}

func TestWithProb(t *testing.T) {
	t.Parallel()

	list := New(WithProb(0.1))
	assert.Equal(t, 0.1, list.prob)
}

func TestWithRandSource(t *testing.T) {
	t.Parallel()

	list := New(WithRandSource(rand.NewSource(2)))
	assert.Equal(t, int64(1543039099823358511), list.randSource.Int63())
}

func TestMakeProbs(t *testing.T) {
	list := &SkipList{
		maxLevel: 4,
		prob:     defaultProb,
	}
	list.makeProbs()
	assert.Len(t, list.probs, 4)

	for i := 1; i < 4; i++ {
		expected := math.Pow(list.prob, float64(i))
		assert.Equal(t, expected, list.probs[i])
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	list := New()
	assert.Nil(t, list.Search(1))

	list.Insert(1, 1)
	assert.Equal(t, 1, list.Search(1))

	list.Insert(2, 2)
	assert.Equal(t, 1, list.Search(1))
	assert.Equal(t, 2, list.Search(2))

	assert.Nil(t, list.Search(3))
}

func TestInsert(t *testing.T) {
	t.Parallel()

	// randLevel will get 3 2 5
	list := New(WithRandSource(rand.NewSource(2)))

	list.Insert(1, 1)
	n1 := list.head.next[0]
	assert.Equal(t, 3, len(n1.next))
	assert.Equal(t, 3, list.level)
	assert.Equal(t, 1, list.Len())

	list.Insert(3, 3)
	n3 := list.head.next[0].next[0]
	assert.Equal(t, 2, len(n3.next))
	assert.Equal(t, 3, list.level)
	assert.Equal(t, 2, list.Len())

	list.Insert(2, 2)
	n2 := list.head.next[0].next[0]
	assert.Equal(t, 5, len(n2.next))
	assert.Equal(t, 5, list.level)
	assert.Equal(t, 3, list.Len())

	list.Insert(3, 4)
	n3 = list.head.next[0].next[0].next[0]
	assert.Equal(t, 4, n3.value)
	assert.Equal(t, 3, list.Len())
}

func TestDelete(t *testing.T) {
	t.Parallel()

	// randLevel will get 3 2 5
	list := New(WithRandSource(rand.NewSource(2)))
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	assert.Equal(t, 5, list.level)
	assert.Equal(t, 3, list.Len())

	list.Delete(1)
	assert.Equal(t, 5, list.level)
	assert.Equal(t, 2, list.Len())

	list.Delete(3)
	assert.Equal(t, 2, list.level)
	assert.Equal(t, 1, list.Len())

	list.Delete(3)
	assert.Equal(t, 2, list.level)
	assert.Equal(t, 1, list.Len())

	list.Delete(2)
	assert.Equal(t, 1, list.level)
	assert.Equal(t, 0, list.Len())
}

func TestPop(t *testing.T) {
	t.Parallel()

	list := New()
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	assert.Equal(t, 1, list.Pop(1))
	assert.Equal(t, 2, list.Pop(2))
	assert.Equal(t, 3, list.Pop(3))
	assert.Equal(t, nil, list.Pop(4))
}

func TestString(t *testing.T) {
	t.Parallel()

	// randLevel will get 3 2 5
	list := New(WithRandSource(rand.NewSource(2)))

	var lines []string

	list.Insert(1, 1)
	lines = strings.Split(list.String(), "\n")
	assert.Len(t, lines, 4)
	assert.Contains(t, lines[1], "(1)")
	assert.Contains(t, lines[2], "(1)")

	list.Insert(2, 2)
	lines = strings.Split(list.String(), "\n")
	assert.Len(t, lines, 4)
	assert.Contains(t, lines[1], "(2)")

	list.Insert(3, 3)
	lines = strings.Split(list.String(), "\n")
	assert.Len(t, lines, 6)
	assert.Contains(t, lines[1], "(3)")
	assert.Contains(t, lines[2], "(3)")
	assert.Contains(t, lines[3], "(3)")
	assert.Contains(t, lines[4], "(3)")
}

var searchResult interface{} = nil

// go test -v -run=^$ -bench=BenchmarkSearch -benchmem -count=10
func BenchmarkSearch(b *testing.B) {
	list := getList()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchResult = list.Search(5000)
	}

	if searchResult.(int) != 5000 {
		b.Errorf("want 50, got %v", searchResult)
	}
}

// go test -v -run=^$ -bench=BenchmarkInsertAndDelete -benchmem -count=10
func BenchmarkInsertAndDelete(b *testing.B) {
	list := getList()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		list.Delete(5000)
		list.Insert(5000, 5000)
	}
}

func getList() *SkipList {
	list := New(WithRandSource(rand.NewSource(2)))
	for i := 0; i < 10000; i++ {
		list.Insert(float64(i), i)
	}

	return list
}
