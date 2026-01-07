// Package testplus includes utilities to ease writing tests.
// Functionality in this package is intended only for test cases and focuses on usability over performance.
//
// All functionality is concurrency-safe.
package testplus

import (
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"sync"
)

//#region ease-of-use globals

// internal, package-level PG
var glbPG *PortGenerator = NewPortGenerator()

// UniquePort returns a new, unique port between [1024-65535].
//
// Panics if more than (65535-1024) unique ports are requested.
func UniquePort() uint16 {
	return glbPG.Generate()
}

//#endregion

// A PortGenerator TODO
//
// PortGenerators are safe for concurrent use.
// Each PortGenerator has its own uniqueness requirement and they are not enforced across instances.
type PortGenerator struct {
	mu sync.Mutex
	//availableWK []uint16 // list of available "well-known" ports
	available []int // list of available ports > 1023
}

// NewPortGenerator returns a struct capable of producing unique port numbers.
// Front-loads the performance cost by generating a full list of valid port numbers to be drawn from later.
func NewPortGenerator() *PortGenerator {
	pg := &PortGenerator{
		available: rand.Perm(math.MaxUint16 - 1023),
	}

	return pg
}

// Generate returns a new, unique port between [1024-65535].
//
// Panics if more than (65535-1024) unique ports are requested.
func (gen *PortGenerator) Generate() uint16 {
	gen.mu.Lock()
	defer gen.mu.Unlock()

	if len(gen.available) == 0 {
		panic("out of port numbers")
	}
	// dequeue
	n := gen.available[0]
	gen.available = gen.available[1:len(gen.available)]

	return uint16(n)
}

// SlicesUnorderedEqual compares the elements of the given slices for equality and equal count without taking order of the elements into account.
func SlicesUnorderedEqual[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	// convert each slice into a map of key --> count
	var wg sync.WaitGroup

	am := make(map[T]uint)
	wg.Go(func() {
		for _, k := range a {
			am[k] += 1
		}
	})

	bm := make(map[T]uint)
	wg.Go(func() {
		for _, k := range b {
			bm[k] += 1
		}
	})

	wg.Wait()
	return maps.Equal(am, bm)
}

// ExpectedActual returns a string comparing the expected result to the actual result over two lines.
// Used to add clarity to unit test error messages.
func ExpectedActual[T any](expected, actual T) string {
	return fmt.Sprintf("\tExpected: '%v'\n\tActual: '%v'", expected, actual)
}
