package cacher

import (
	"sync"
)

type Map struct {
	m   sync.Map
	gen func(key interface{}) Cacher
}

func NewMap(cacherGen func(key interface{}) Cacher) *Map {
	return &Map{
		m:   sync.Map{},
		gen: cacherGen,
	}
}

func (m *Map) Get(key interface{}) Cacher {
	c, ok := m.m.Load(key)
	if !ok {
		c, _ = m.m.LoadOrStore(key, m.gen(key))
	}
	return c.(Cacher)
}

func (m *Map) Delete(key interface{}) {
	m.m.Delete(key)
}
