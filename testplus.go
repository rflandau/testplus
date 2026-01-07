// Package testplus includes utilities to ease writing tests.
// Functionality in this package is intended only for test cases and focuses on usability over performance.
package testplus

import (
	"maps"
	"math"
	"math/rand/v2"
	"sync"
)

//#region ease-of-use globals

// internal, package-level PG
var glbPG *PortGenerator = NewPortGenerator()

// UniquePort generates a package-level-unique port.
// Will only return a port < 1024 if includeWellKnown.
//
// Uniqueness is NOT enforced relative to local PortGenerators.
/*func UniquePort(includeWellKnown bool) uint16 {
	return glbPG.Generate(includeWellKnown)
}*/

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
// It front-loads the performance cost by generating a full list of valid port numbers to be drawn from later
func NewPortGenerator() *PortGenerator {
	pg := &PortGenerator{
		//availableWK: make([]uint16, 1023),                // (1-1023)
		available: rand.Perm(math.MaxUint16 - 1023),
	}

	return pg
}

// Generate returns a new, unique port between [1024-65535].
//
// If no ports are available... TODO
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
