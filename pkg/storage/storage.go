package storage

import (
	"sync"
	"time"
)

type Storage struct {
	sync.RWMutex
	data storageMap
}

type storageData struct {
	Value any
	TTL   time.Duration
}

type storageMap map[string]storageData

func (s *Storage) Init() {
	sMap := storageMap{}
	s.data = sMap
}

func (s *Storage) Add(key string, value any, ttl time.Duration) {
	sData := storageData{
		Value: value,
		TTL:   ttl,
	}
	s.Lock()
	s.data[key] = sData
	s.Unlock()
	go s.clearData(key, ttl)
}

func (s *Storage) Get(key string) storageData {
	s.RLock()
	s.RUnlock()
	return s.data[key]
}

func (s *Storage) clearData(key string, ttl time.Duration) {
	time.Sleep(ttl)
	delete(s.data, key)
}
