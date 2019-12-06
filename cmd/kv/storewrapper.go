package main

import (
	"fmt"

	"github.com/abronan/valkeyrie/store"
)

type storeWrapper struct {
	store.Store
}

func (s *storeWrapper) Put(key string, value []byte, options *store.WriteOptions) error {
	fmt.Println("Put", key, string(value))

	if s.Store == nil {
		return nil
	}
	return s.Store.Put(key, value, options)
}

func (s *storeWrapper) Get(key string, options *store.ReadOptions) (*store.KVPair, error) {
	fmt.Println("Get", key)

	if s.Store == nil {
		return nil, nil
	}
	return s.Store.Get(key, options)
}

func (s *storeWrapper) Delete(key string) error {
	fmt.Println("Delete", key)

	if s.Store == nil {
		return nil
	}
	return s.Store.Delete(key)
}

func (s *storeWrapper) Exists(key string, options *store.ReadOptions) (bool, error) {
	fmt.Println("Exists", key)

	if s.Store == nil {
		return true, nil
	}
	return s.Store.Exists(key, options)
}

func (s *storeWrapper) Watch(key string, stopCh <-chan struct{}, options *store.ReadOptions) (<-chan *store.KVPair, error) {
	fmt.Println("Watch", key)

	if s.Store == nil {
		return nil, nil
	}
	return s.Store.Watch(key, stopCh, options)
}

func (s *storeWrapper) WatchTree(directory string, stopCh <-chan struct{}, options *store.ReadOptions) (<-chan []*store.KVPair, error) {
	fmt.Println("WatchTree", directory)

	if s.Store == nil {
		return nil, nil
	}
	return s.Store.WatchTree(directory, stopCh, options)
}

func (s *storeWrapper) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	fmt.Println("NewLock", key)

	if s.Store == nil {
		return nil, nil
	}
	return s.Store.NewLock(key, options)
}

func (s *storeWrapper) List(directory string, options *store.ReadOptions) ([]*store.KVPair, error) {
	fmt.Println("List", directory)

	if s.Store == nil {
		return nil, nil
	}
	return s.Store.List(directory, options)
}

func (s *storeWrapper) DeleteTree(directory string) error {
	fmt.Println("DeleteTree", directory)

	if s.Store == nil {
		return nil
	}
	return s.Store.DeleteTree(directory)
}

func (s *storeWrapper) AtomicPut(key string, value []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) {
	fmt.Println("AtomicPut", key, string(value), previous)

	if s.Store == nil {
		return true, nil, nil
	}
	return s.Store.AtomicPut(key, value, previous, options)
}

func (s *storeWrapper) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	fmt.Println("AtomicDelete", key, previous)

	if s.Store == nil {
		return true, nil
	}
	return s.Store.AtomicDelete(key, previous)
}

func (s *storeWrapper) Close() {
	if s.Store == nil {
		return
	}
	s.Store.Close()
}
