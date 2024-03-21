package memstore

import (
	"sync"
)

type Item struct {
	ID   string
	Name string
	Lat  float64
	Lon  float64
}

func New() *Store {
	return &Store{items: make([]Item, 0)}
}

type Store struct {
	mu    *sync.RWMutex
	items []Item
}

func (s *Store) FindAll() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Item, len(s.items))
	copy(items, s.items)
	return items
}

func (s *Store) Add(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(
		s.items, Item{
			ID:   item.ID,
			Name: item.Name,
			Lat:  item.Lat,
			Lon:  item.Lon,
		},
	)
}

func (s *Store) DeleteByObjectID(objectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for idx, item := range s.items {
		if item.ID == objectID {
			s.items = append(s.items[:idx], s.items[idx+1:]...)
			return
		}
	}
}

func (s *Store) FindByObjectID(objectID string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.ID == objectID {
			return item, true
		}
	}
	return Item{}, false
}
