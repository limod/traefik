package main

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	"github.com/abronan/valkeyrie/store/dynamodb"
	etcdv3 "github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/abronan/valkeyrie/store/redis"
	"github.com/abronan/valkeyrie/store/zookeeper"
)

func main() {
	targets := []store.Backend{store.CONSUL, store.ETCDV3, store.REDIS, store.ZK}

	for _, target := range targets {
		kvStore, err := createClientStore(target)
		if err != nil {
			log.Fatal(err)
		}

		dynConfPath := "./cmd/kv/file.toml"
		conf := map[string]interface{}{}
		_, err = toml.DecodeFile(dynConfPath, &conf)
		if err != nil {
			log.Fatal(err)
		}

		c := client{store: kvStore}
		err = c.load("traefik", conf)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func createClientStore(backend store.Backend) (store.Store, error) {
	switch backend {
	case store.CONSUL:
		addrs := []string{"localhost:8500"}
		return createStore(backend, addrs)

	case store.ETCDV3:
		addrs := []string{"localhost:2379"}
		return createStore(backend, addrs)

	case store.ZK:
		addrs := []string{"localhost:2181"}
		return createStore(backend, addrs)

	case store.REDIS:
		addrs := []string{"localhost:6379"}
		return createStore(backend, addrs)

	case store.DYNAMODB:
		addrs := []string{"http://localhost:8000"}
		return createStore(backend, addrs)

	default:
		return &storeWrapper{}, nil
	}
}

func createStore(backend store.Backend, addrs []string) (store.Store, error) {
	storeConfig := &store.Config{
		ConnectionTimeout: 3 * time.Second,
		Bucket:            "traefik",
	}

	switch backend {
	case store.CONSUL:
		consul.Register()
	case store.ETCDV3:
		etcdv3.Register()
	case store.ZK:
		zookeeper.Register()
	case store.BOLTDB:
		boltdb.Register()
	case store.REDIS:
		redis.Register()
	case store.DYNAMODB:
		dynamodb.Register()
	}

	kvStore, err := valkeyrie.NewStore(backend, addrs, storeConfig)
	if err != nil {
		return nil, err
	}

	return &storeWrapper{Store: kvStore}, nil
}

type client struct {
	store store.Store
}

func (c client) load(parentKey string, conf map[string]interface{}) error {
	for k, v := range conf {
		switch entry := v.(type) {
		case map[string]interface{}:
			key := path.Join(parentKey, k)

			if len(entry) == 0 {
				err := c.store.Put(key, nil, nil)
				if err != nil {
					return err
				}
			} else {
				if err := c.load(key, entry); err != nil {
					return err
				}
			}
		case []map[string]interface{}:
			for i, o := range entry {
				key := path.Join(parentKey, k, strconv.Itoa(i))

				if err := c.load(key, o); err != nil {
					return err
				}
			}
		case []interface{}:
			for i, o := range entry {
				key := path.Join(parentKey, k, strconv.Itoa(i))

				err := c.store.Put(key, []byte(fmt.Sprintf("%v", o)), nil)
				if err != nil {
					return err
				}
			}
		default:
			key := path.Join(parentKey, k)

			err := c.store.Put(key, []byte(fmt.Sprintf("%v", v)), nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
