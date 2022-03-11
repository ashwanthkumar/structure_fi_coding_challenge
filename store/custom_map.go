package store

import (
	"github.com/dgryski/go-mph"
)

// MapWithPHF is a Hash map implementation that uses a Minimal Perfect Hash Function that is discovered based on the known set of keys.
// Good introduction about the approach the related research is described at: https://blog.gopheracademy.com/advent-2017/mphf/
type MapWithPHF struct {
	Length int
	Keys   []string
	Values []*Stat

	hashFunction *mph.Table
}

func NewMapWithPHF(keys []string) *MapWithPHF {
	mph := mph.New(keys)

	return &MapWithPHF{
		Length:       len(keys),
		Keys:         keys,
		hashFunction: mph,
		Values:       make([]*Stat, len(keys)),
	}
}

// Idea behind using go-mph is to identify a location that is unique and mod that value to the total length that we have
// so we can map that to our finite array than always maintain an array of length int32. Since the MPH is built specifically
// for the known set of keys, we should not have any collisions.
func (m *MapWithPHF) Set(key string, value Stat) {
	index := m.hashFunction.Query(key)
	location := index % int32(m.Length)
	m.Values[location] = &value
}

func (m *MapWithPHF) Get(key string) (*Stat, bool) {
	index := m.hashFunction.Query(key)
	location := index % int32(m.Length)
	value := m.Values[location]
	if value != nil {
		return value, true
	} else {
		return nil, false
	}
}
