package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/containous/traefik/configuration"
	"github.com/containous/traefik/types"
)

// /home/ldez/sources/go/bin/gometalinter --vendor ./...

func TestDynConfiguration(t *testing.T) {

	config := &types.Configuration{}

	hydration := NewHydration()
	hydration.Fill(&config)

	f, _ := os.Create("./config-dyn.toml")
	toml.NewEncoder(f).Encode(config)
	defer func() { _ = f.Close() }()

	f, _ = os.Create("./config-dyn.json")
	json.NewEncoder(f).Encode(config)
	defer func() { _ = f.Close() }()
}

func TestEntryPoint(t *testing.T) {

	entryPoints := configuration.EntryPoints{}

	hydration := NewHydration()
	hydration.Fill(&entryPoints)

	f, _ := os.Create("./config-entryPoint.toml")
	toml.NewEncoder(f).Encode(entryPoints)
	defer func() { _ = f.Close() }()
}

func TestGlobalConfiguration(t *testing.T) {
	config := configuration.GlobalConfiguration{}
	//config := tls.Config{}

	hydration := NewHydration()
	//hydration.ExcludedFieldNames = []string{"ACME", "Stats", "StatsRecorder", "CurrentConfigurations"}
	hydration.ExcludedFieldNames = []string{"ACME.TLSConfig", "Stats", "StatsRecorder", "CurrentConfigurations", "API.DashboardAssets"}
	//hydration.ExcludedFieldNames = []string{"Stats", "StatsRecorder", "CurrentConfigurations"}

	hydration.Fill(&config)

	generateKVContent(&config)

	fmt.Printf("%+v\n", config.ACME)

	f, _ := os.Create("./config-global.toml")
	toml.NewEncoder(f).Encode(config)
	defer func() { _ = f.Close() }()

	f, _ = os.Create("./config-global.json")
	json.NewEncoder(f).Encode(config)
	defer func() { _ = f.Close() }()
}

func TestGen(t *testing.T) {
	hu := struct {
		Foo struct {
			One string
			Two string
		}
		Bar string
		Fii []string
		Bir int64
		Fuu map[string]string
		Bur bool
		Faa map[string]int
		Byr map[string]struct {
			One string
			Two string
		}
		Bor map[string][]string
	}{}

	hydration := NewHydration()
	hydration.Fill(&hu)

	fmt.Printf("%+v\n", hu)

	f, _ := os.Create("./test.toml")
	toml.NewEncoder(f).Encode(hu)
}

func TestName(t *testing.T) {

	config := &types.Configuration{}
	if _, err := toml.DecodeFile("./config-dyn-test.toml", config); err != nil {
		fmt.Println(err)
		return
	}
}
