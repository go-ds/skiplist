package skiplist

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
	list := New(WithRandSource(rand.NewSource(2)))
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	fmt.Println(list.String())

	// output:
	// level  1 --> 1.000000(1) <--> 2.000000(2) <--> 3.000000(3) --> nil
	// level  2 --> 1.000000(1) <--> 2.000000(2) <--> 3.000000(3) --> nil
	// level  3 --> 3.000000(3) --> nil
}

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

func TestDisableMutex(t *testing.T) {
	t.Parallel()

	list := New(DisableMutex())
	assert.Equal(t, true, list.disableMutex)
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

func ExampleSkipList_Search() {
	list := New()
	fmt.Println(list.Search(1))

	list.Insert(1, 1)
	fmt.Println(list.Search(1))

	list.Insert(2, 2)
	fmt.Println(list.Search(1))
	fmt.Println(list.Search(2))
	fmt.Println(list.Search(3))
	// output:
	// <nil>
	// 1
	// 1
	// 2
	// <nil>
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
	list := New(WithRandSource(rand.NewSource(2)), WithProb(0.5))

	list.Insert(1, 1)
	n1 := list.head.next[0]
	assert.Equal(t, 3, len(n1.next))
	assert.Equal(t, 3, list.level)
	assert.Equal(t, 1, list.Size())

	list.Insert(3, 3)
	n3 := list.head.next[0].next[0]
	assert.Equal(t, 2, len(n3.next))
	assert.Equal(t, 3, list.level)
	assert.Equal(t, 2, list.Size())

	list.Insert(2, 2)
	n2 := list.head.next[0].next[0]
	assert.Equal(t, 5, len(n2.next))
	assert.Equal(t, 5, list.level)
	assert.Equal(t, 3, list.Size())

	list.Insert(3, 4)
	n3 = list.head.next[0].next[0].next[0]
	assert.Equal(t, 4, n3.value)
	assert.Equal(t, 3, list.Size())
}

func ExampleSkipList_Delete() {
	list := New(WithRandSource(rand.NewSource(2)))
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	list.Delete(2)

	fmt.Println(list.String())

	// output:
	// level  1 --> 1.000000(1) <--> 3.000000(3) --> nil
	// level  2 --> 1.000000(1) <--> 3.000000(3) --> nil
	// level  3 --> 3.000000(3) --> nil
}

func TestDelete(t *testing.T) {
	t.Parallel()

	// randLevel will get 3 2 5
	list := New(WithRandSource(rand.NewSource(2)), WithProb(0.5))
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	assert.Equal(t, 5, list.level)
	assert.Equal(t, 3, list.Size())

	list.Delete(1)
	assert.Equal(t, 5, list.level)
	assert.Equal(t, 2, list.Size())

	list.Delete(3)
	assert.Equal(t, 2, list.level)
	assert.Equal(t, 1, list.Size())

	list.Delete(3)
	assert.Equal(t, 2, list.level)
	assert.Equal(t, 1, list.Size())

	list.Delete(2)
	assert.Equal(t, 1, list.level)
	assert.Equal(t, 0, list.Size())
}

func ExampleSkipList_Pop() {
	list := New(WithRandSource(rand.NewSource(2)))
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	fmt.Println(list.Pop(2))
	fmt.Println(list.Pop(4))

	fmt.Println(list.String())

	// output:
	// 2
	// <nil>
	// level  1 --> 1.000000(1) <--> 3.000000(3) --> nil
	// level  2 --> 1.000000(1) <--> 3.000000(3) --> nil
	// level  3 --> 3.000000(3) --> nil
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

func ExampleSkipList_Clear() {
	list := New(WithRandSource(rand.NewSource(2)))
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	list.Clear()

	fmt.Println(list.String())

	// output:
	// level  1 --> nil
}

func TestClear(t *testing.T) {
	t.Parallel()

	list := New()
	list.Insert(1, 1)
	list.Insert(2, 2)
	list.Insert(3, 3)

	list.Clear()
	assert.Equal(t, 0, list.Size())
	assert.Equal(t, true, list.Empty())
	assert.Equal(t, "level  1 --> nil\n", list.String())
}

func TestString(t *testing.T) {
	t.Parallel()

	// randLevel will get 3 2 5
	list := New(WithRandSource(rand.NewSource(2)), WithProb(0.5))

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

// go test -v -run=^$ -bench=BenchmarkSearch -benchmem -count=4
func BenchmarkSearch1000(b *testing.B) {
	benchmarkSearch(b, 1000)
}

func BenchmarkSearch10000(b *testing.B) {
	benchmarkSearch(b, 10000)
}

func BenchmarkSearch100000(b *testing.B) {
	benchmarkSearch(b, 100000)
}

func BenchmarkSearch1000000(b *testing.B) {
	benchmarkSearch(b, 1000000)
}

func benchmarkSearch(b *testing.B, n int) {
	list := New(WithRandSource(rand.NewSource(2)))
	for i := 0; i < n; i++ {
		list.Insert(float64(n-i), i)
	}
	target := float64(n / 2)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			searchResult = list.Search(target)
		}
	})

	assert.Equal(b, n/2, searchResult)
}

// go test -v -run=^$ -bench=BenchmarkInsertAndDelete -benchmem -count=4
func BenchmarkInsertAndDelete100(b *testing.B) {
	benchmarkInsertAndDelete(b, 100)
}

func BenchmarkInsertAndDelete1000(b *testing.B) {
	benchmarkInsertAndDelete(b, 1000)
}

func BenchmarkInsertAndDelete10000(b *testing.B) {
	benchmarkInsertAndDelete(b, 10000)
}

func BenchmarkInsertAndDelete100000(b *testing.B) {
	benchmarkInsertAndDelete(b, 100000)
}

func BenchmarkInsertAndDelete1000000(b *testing.B) {
	benchmarkInsertAndDelete(b, 1000000)
}

func benchmarkInsertAndDelete(b *testing.B, n int) {
	list := New(WithRandSource(rand.NewSource(2)))
	for i := 0; i < n; i++ {
		list.Insert(float64(n-i), i)
	}
	target := float64(n / 2)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list.Delete(target)
			list.Insert(target, n/2)
		}
	})
}

// go test -v -run=^$ -bench=BenchmarkBestInsert -benchmem -count=4
func BenchmarkBestInsert100(b *testing.B) {
	benchmarkBestInsert(b, 100)
}

func BenchmarkBestInsert1000(b *testing.B) {
	benchmarkBestInsert(b, 1000)
}

func BenchmarkBestInsert10000(b *testing.B) {
	benchmarkBestInsert(b, 10000)
}

func BenchmarkBestInsert100000(b *testing.B) {
	benchmarkBestInsert(b, 100000)
}

func BenchmarkBestInsert1000000(b *testing.B) {
	benchmarkBestInsert(b, 1000000)
}

func benchmarkBestInsert(b *testing.B, n int) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list := New(WithRandSource(rand.NewSource(2)))
			for j := 0; j < n; j++ {
				list.Insert(float64(n-j), j)
			}
		}
	})
}

// go test -v -run=^$ -bench=BenchmarkWorstInsert -benchmem -count=4
func BenchmarkWorstInsert100(b *testing.B) {
	benchmarkWorstInsert(b, 100)
}

func BenchmarkWorstInsert1000(b *testing.B) {
	benchmarkWorstInsert(b, 1000)
}

func BenchmarkWorstInsert10000(b *testing.B) {
	benchmarkWorstInsert(b, 10000)
}

func BenchmarkWorstInsert100000(b *testing.B) {
	benchmarkWorstInsert(b, 100000)
}

func BenchmarkWorstInsert1000000(b *testing.B) {
	benchmarkWorstInsert(b, 1000000)
}

func benchmarkWorstInsert(b *testing.B, n int) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list := New()
			for j := 0; j < n; j++ {
				list.Insert(float64(j), j)
			}
		}
	})
}

// go test -v -run=^$ -bench=BenchmarkBestDelete -benchmem -count=4
func BenchmarkBestDelete100(b *testing.B) {
	benchmarkBestDelete(b, 100)
}

func BenchmarkBestDelete1000(b *testing.B) {
	benchmarkBestDelete(b, 1000)
}

func BenchmarkBestDelete10000(b *testing.B) {
	benchmarkBestDelete(b, 10000)
}

func BenchmarkBestDelete100000(b *testing.B) {
	benchmarkBestDelete(b, 100000)
}

func BenchmarkBestDelete1000000(b *testing.B) {
	benchmarkBestDelete(b, 1000000)
}

func benchmarkBestDelete(b *testing.B, n int) {
	list := New(WithRandSource(rand.NewSource(2)))
	for i := 0; i < n; i++ {
		list.Insert(float64(n-i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list.Delete(0)
		}
	})
}

// go test -v -run=^$ -bench=BenchmarkWorstDelete -benchmem -count=4
func BenchmarkWorstDelete100(b *testing.B) {
	benchmarkWorstDelete(b, 100)
}

func BenchmarkWorstDelete1000(b *testing.B) {
	benchmarkWorstDelete(b, 1000)
}

func BenchmarkWorstDelete10000(b *testing.B) {
	benchmarkWorstDelete(b, 10000)
}

func BenchmarkWorstDelete100000(b *testing.B) {
	benchmarkWorstDelete(b, 100000)
}

func BenchmarkWorstDelete1000000(b *testing.B) {
	benchmarkWorstDelete(b, 1000000)
}

func benchmarkWorstDelete(b *testing.B, n int) {
	list := New(WithRandSource(rand.NewSource(2)))
	for i := 0; i < n; i++ {
		list.Insert(float64(n-i), i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			list.Delete(float64(n - 1))
		}
	})
}
