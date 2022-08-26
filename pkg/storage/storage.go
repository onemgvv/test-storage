package storage

import (
	"time"
)

type Storage struct {
	data StorageMap
}

type StorageData struct {
	Value any
	TTL   time.Duration
}

type StorageMap map[string]StorageData

func (s *Storage) Init() {
	sMap := StorageMap{}
	s.data = sMap
}

func (s *Storage) Add(key string, value any, ttl time.Duration) {
	sData := StorageData{
		Value: value,
		TTL: ttl,
	}

	s.data[key] = sData
	go s.clearData(key, ttl)
}

func (s *Storage) Get(key string) StorageData {
	return s.data[key]
}

func (s *Storage) clearData(key string, ttl time.Duration) {
	time.Sleep(ttl)
	delete(s.data, key)
}