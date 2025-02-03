package storage

import (
	"errors"

	"github.com/aswinbennyofficial/raikuv/internal/storage/hmap"
	"github.com/aswinbennyofficial/raikuv/internal/storage/lsm"
)

const (
	InMemory string = "map"
	LSMTree  string = "lsm"
)

type Storage interface {
    Get(key string) (interface{}, bool)
    Put(key string, value interface{})
    Pop(key string)
}

func NewStorage(Storagetype string)(Storage,error){
	switch Storagetype{
	case InMemory :
		return hmap.NewMapStore(),nil
	case LSMTree :
		return lsm.NewLSMStore(),nil
	default:
		return nil, errors.New("unsupported storage backend")
	}
	
}


