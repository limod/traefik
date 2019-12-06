package main

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/abronan/valkeyrie/store"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/kv"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	targets := []store.Backend{store.CONSUL, store.ETCDV3, store.REDIS, store.ZK}

	for _, target := range targets {
		kvStore, err := createClientStore(target)
		require.NoError(t, err)

		root := "traefik"

		pairs, err := kvStore.List(root, nil)
		require.NoError(t, err)

		cfg := &dynamic.Configuration{}
		err = kv.Decode(pairs, cfg, root)
		require.NoError(t, err)

		file, err := os.Create("./" + string(target) + ".toml")
		require.NoError(t, err)

		err = toml.NewEncoder(file).Encode(cfg)
		require.NoError(t, err)

		// values := squash(pairs)
		// fmt.Println()
		// for k, v := range values {
		// 	fmt.Println(k, ":", v)
		// }
		//
		// err = parser.Decode(values, cfg, "traefik", "traefik.http")
		// require.NoError(t, err)
	}
}

func squash(pairs []*store.KVPair) map[string]string {
	exp := regexp.MustCompile(`^(.+)/\d+$`)

	values := map[string]string{}
	for _, pair := range pairs {
		// fmt.Println(pair.Key, pair.Value)

		if exp.MatchString(pair.Key) {
			sanitizedKey := toLabel(exp.FindStringSubmatch(pair.Key)[1])

			if v, ok := values[sanitizedKey]; ok {
				values[sanitizedKey] = v + "," + string(pair.Value)
			} else {
				values[sanitizedKey] = string(pair.Value)
			}
		} else {
			sanitizedKey := toLabel(pair.Key)

			if v, ok := values[sanitizedKey]; ok {
				if len(pair.Value) != 0 && v != "" {
					values[sanitizedKey] = v + "," + string(pair.Value)
				}
			} else {
				values[sanitizedKey] = string(pair.Value)
			}
		}
	}

	return values
}

func toLabel(name string) string {
	fargs := func(c rune) bool {
		return c == '/'
	}

	return strings.Join(strings.FieldsFunc(name, fargs), ".")
}
