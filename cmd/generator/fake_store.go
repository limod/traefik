package generator

import (
	"fmt"

	"github.com/abronan/valkeyrie/store"
	"github.com/containous/staert"
)

func generateKVContent(data interface{}) error {
	kv := &staert.KvSource{
		Store:  &FakeStore{},
		Prefix: "traefik",
	}
	return kv.StoreConfig(data)
}

type FakeStore struct{}

func (FakeStore) Put(key string, value []byte, options *store.WriteOptions) error {
	if len(value) > 0 {
		fmt.Printf("| `%s` | `%s` |\n", key, string(value))
	} else {
		fmt.Printf("| `%s` |   |\n", key)
	}

	return nil
}

func (FakeStore) Get(key string, options *store.ReadOptions) (*store.KVPair, error) {
	panic("implement me")
}

func (FakeStore) Delete(key string) error {
	panic("implement me")
}

func (FakeStore) Exists(key string, options *store.ReadOptions) (bool, error) {
	panic("implement me")
}

func (FakeStore) Watch(key string, stopCh <-chan struct{}, options *store.ReadOptions) (<-chan *store.KVPair, error) {
	panic("implement me")
}

func (FakeStore) WatchTree(directory string, stopCh <-chan struct{}, options *store.ReadOptions) (<-chan []*store.KVPair, error) {
	panic("implement me")
}

func (FakeStore) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	panic("implement me")
}

func (FakeStore) List(directory string, options *store.ReadOptions) ([]*store.KVPair, error) {
	panic("implement me")
}

func (FakeStore) DeleteTree(directory string) error {
	panic("implement me")
}

func (FakeStore) AtomicPut(key string, value []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) {
	panic("implement me")
}

func (FakeStore) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	panic("implement me")
}

func (FakeStore) Close() {
	panic("implement me")
}
