package convert

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestStaticTomlToYaml(t *testing.T) {
	conf := &static.Configuration{}

	// generator.Generate(conf)

	// _, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/static-configuration/file.toml", conf)
	_, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/exp/convert/bla.toml", conf)
	require.NoError(t, err)

	out, err := yaml.Marshal(conf)
	require.NoError(t, err)

	fmt.Println(string(out))
}

func TestStaticTomlToToml(t *testing.T) {
	conf := &static.Configuration{}

	// generator.Generate(conf)

	// _, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/static-configuration/file.toml", conf)
	_, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/exp/convert/bla.toml", conf)
	require.NoError(t, err)

	err = toml.NewEncoder(os.Stdout).Encode(conf)
	require.NoError(t, err)
}

func TestStaticYamlToToml(t *testing.T) {
	conf := &static.Configuration{}

	// generator.Generate(conf)

	// content, err := ioutil.ReadFile("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/static-configuration/file.yml")
	content, err := ioutil.ReadFile("/home/ldez/sources/go/src/github.com/containous/traefik/exp/convert/bla.yml")
	require.NoError(t, err)

	err = yaml.Unmarshal(content, conf)
	require.NoError(t, err)

	err = toml.NewEncoder(os.Stdout).Encode(conf)
	require.NoError(t, err)
}

func TestDynTomlToYaml(t *testing.T) {
	conf := &dynamic.Configuration{}

	// generator.Generate(conf)

	// _, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/file.yaml", conf)
	_, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/exp/convert/bla.toml", conf)

	out, err := yaml.Marshal(conf)
	require.NoError(t, err)

	fmt.Println(string(out))
}

func TestDynTomlToToml(t *testing.T) {
	conf := &dynamic.Configuration{}

	// generator.Generate(conf)

	_, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/file.toml", conf)
	// _, err := toml.DecodeFile("/home/ldez/sources/go/src/github.com/containous/traefik/exp/convert/bla.toml", conf)
	require.NoError(t, err)

	err = toml.NewEncoder(os.Stdout).Encode(conf)
	require.NoError(t, err)
}

func TestDynYamlToToml(t *testing.T) {
	conf := &dynamic.Configuration{}

	// generator.Generate(conf)

	// content, err := ioutil.ReadFile("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/file.yml")
	content, err := ioutil.ReadFile("/home/ldez/sources/go/src/github.com/containous/traefik/exp/convert/bla.yml")
	require.NoError(t, err)

	err = yaml.Unmarshal(content, conf)
	require.NoError(t, err)

	err = toml.NewEncoder(os.Stdout).Encode(conf)
	require.NoError(t, err)
}
