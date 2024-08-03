package OrderedMap

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"
)

var c OrderMapImp

func init() {
	c = OrderMapImp{}
	c.Init()
}

func TestOrder(t *testing.T) {
	round := 100
	var expected []string
	for i := 0; i < round; i++ {
		expected = append(expected, strconv.Itoa(i))
		c.Add(strconv.Itoa(i), strconv.Itoa(i))
	}

	keys := c.Range()

	for i := 0; i < round; i++ {
		if keys[i] != expected[i] {
			t.Fatalf("case %d failed", i)
			t.Failed()
		}
	}
}

func TestParallel_Add_Range(t *testing.T) {

	concurrency := 100
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < concurrency; i++ {
		go func(idx int, ctx context.Context) {
			index := strconv.Itoa(idx)
			cctx, _ := context.WithCancel(ctx)
			var duration time.Duration
			for duration <= 0 {
				duration = time.Duration(rand.Intn(5))
			}
			ticker := time.NewTicker(time.Second * duration)

			for {
				select {
				case <-ticker.C:
					if idx%2 == 0 {
						keys := c.Range()
						if sort.StringsAreSorted(keys) {
							t.Failed()
						}
					} else {
						c.Add(index, index)
					}
				case <-cctx.Done():
					return
				}
			}

		}(i, ctx)
	}
	time.Sleep(time.Second * 10)
	cancel()
}

func TestParallel_Del_Range(t *testing.T) {
	concurrency := 100
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < concurrency; i++ {
		idx := i
		c.Add(strconv.Itoa(idx), strconv.Itoa(idx))
	}
	for i := 0; i < concurrency; i++ {
		idx := i
		go func(idx int, ctx context.Context) {
			index := strconv.Itoa(idx)
			cctx, _ := context.WithCancel(ctx)
			var duration time.Duration
			for duration <= 0 {
				duration = time.Duration(rand.Intn(5))
			}
			ticker := time.NewTicker(time.Second * duration)

			if idx%2 == 0 {
				for {
					select {
					case <-ticker.C:
						keys := c.Range()
						for _, key := range keys {
							numkey, err := strconv.Atoi(key)
							if err != nil {
								t.Fatal(err)
							}
							if numkey%2 == 1 {
								t.Fatalf("case %d failed", numkey)
								t.Failed()
							}
						}
					case <-cctx.Done():
						return
					}
				}
			} else {
				c.Del(index)
			}

		}(idx, ctx)
	}
	time.Sleep(time.Second * 10)
	cancel()
}
