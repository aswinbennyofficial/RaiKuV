package storage

import (	
	"sync"
)



type DataStore struct {
	store sync.Map
};

func NewDataStore() *DataStore{
	return &DataStore{}
} 


// Get Retrieves value from the map
func (dataStore *DataStore) Get(key string)(interface{}, bool){
	value, ok := dataStore.store.Load(key)

	return value,ok
}


func (dataStore *DataStore) Put(key string, value interface{}){
	dataStore.store.Store(key, value)
}

func (dataStore *DataStore) Pop(key string){
	dataStore.store.Delete(key)
}