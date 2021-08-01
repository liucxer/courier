package snowflakeid

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Benchmark(b *testing.B) {
	startTime, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")

	for _, sets := range [][3]uint{
		{10, 12, 1},
		{16, 10, 1},
		{16, 8, 1},
		{16, 8, 5},
		{16, 8, 10},
	} {
		f := NewSnowflakeFactory(sets[0], sets[1], sets[2], startTime)
		g, _ := f.NewSnowflake(1)

		b.Run(fmt.Sprintf("end at %s, max worker %d max sequence %d per %dms", f.MaxTime(), f.MaxWorkerID(), f.MaxSequence(), sets[2]), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := g.ID()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestSnowflake(t *testing.T) {
	suite := NewTestSuite(t, 1000000)
	generator, err := NewSnowflake(1)
	if err != nil {
		t.Fatal(err)
	}

	suite.RunGenerator(generator)
	suite.ExpectN(suite.N)
}

func TestSnowflake_InMultiGoroutine(t *testing.T) {
	suite := NewTestSuite(t, 100)
	generator, err := NewSnowflake(1)
	if err != nil {
		t.Fatal(err)
	}

	n := 1000

	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			suite.RunGenerator(generator)
		}()
	}

	wg.Wait()

	suite.ExpectN(suite.N * n)
}

func TestSnowflake_WithDifferentWorkerID(t *testing.T) {
	suite := NewTestSuite(t, 100000)

	wg := &sync.WaitGroup{}

	n := 10

	for i := 0; i < n; i++ {
		generator, err := NewSnowflake(uint32(i + 1))
		if err != nil {
			t.Fatal(err)
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			suite.RunGenerator(generator)
		}()
	}

	wg.Wait()

	suite.ExpectN(suite.N * n)
}

func NewTestSuite(t *testing.T, n int) *Suite {
	return &Suite{
		T: t,
		N: n,
	}
}

type Suite struct {
	*testing.T
	N int
	sync.Map
}

func (s *Suite) ExpectN(n int) {
	t := s.Total()
	if t != n {
		s.Fatalf("expect generated %d, but go %d", n, t)
	}
}

func (s *Suite) Total() int {
	c := 0
	s.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return c
}

func (s *Suite) RunGenerator(generator *Snowflake) {
	for i := 1; i <= s.N; i++ {
		id, err := generator.ID()
		if err != nil {
			s.Fatal(err)
		}
		s.Store(id, true)
	}
}
