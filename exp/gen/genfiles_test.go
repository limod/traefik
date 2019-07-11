package gen

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/containous/traefik/v2/pkg/api"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/generator"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestGenStaticTOML(t *testing.T) {
	element := &static.Configuration{}

	generator.Generate(element)

	hydration := NewHydration()
	hydration.Fill(element)

	file, err := os.Create(staticFileToml)
	require.NoError(t, err)

	delete(element.EntryPoints, "EntryPoint1")

	err = toml.NewEncoder(file).Encode(element)
	require.NoError(t, err)
}

func TestGenStaticYAML(t *testing.T) {
	element := &static.Configuration{}

	generator.Generate(element)
	hydration := NewHydration()

	hydration.Fill(element)

	file, err := os.Create(staticFileYaml)
	require.NoError(t, err)

	delete(element.EntryPoints, "EntryPoint1")

	out, err := yaml.Marshal(element)
	require.NoError(t, err)

	_, err = file.Write(out)
	require.NoError(t, err)
}

func TestGenDynTOML(t *testing.T) {
	element := &dynamic.Configuration{}
	// generator.Generate(element)

	hydration := NewHydration()
	hydration.Fill(element)

	cleanMiddlewares(element)
	cleanHTTPServices(element)
	cleanTCPServices(element)
	cleanUDPServices(element)

	file, err := os.Create(dynFileToml)
	require.NoError(t, err)

	err = toml.NewEncoder(file).Encode(element)
	require.NoError(t, err)
}

func TestGenDynYAML(t *testing.T) {
	element := &dynamic.Configuration{}
	// generator.Generate(element)

	hydration := NewHydration()
	hydration.Fill(element)

	cleanMiddlewares(element)
	cleanHTTPServices(element)
	cleanTCPServices(element)
	cleanUDPServices(element)

	file, err := os.Create(dynFileYaml)
	require.NoError(t, err)

	out, err := yaml.Marshal(element)
	require.NoError(t, err)

	_, err = file.Write(out)
	require.NoError(t, err)
}

func TestGenAPI(t *testing.T) {
	testCases := []struct {
		desc     string
		element  interface{}
		filename string
	}{
		{
			desc:     "RunTimeRepresentation",
			element:  &api.RunTimeRepresentation{},
			filename: "api_rawdata.json",
		},
		// {
		// 	desc:     "Overview",
		// 	element:  &api.Overview{},
		// 	filename: "api_overview.json",
		// },
		// {
		// 	desc:     "EntryPointRepresentation",
		// 	element:  &api.EntryPointRepresentation{},
		// 	filename: "api_entryPointRepresentation.json",
		// },
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			// generator.Generate(element)

			hydration := NewHydration()
			hydration.Fill(test.element)

			// cleanMiddlewares(element)
			// cleanServers(element)

			file, err := os.Create("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/" + test.filename)
			require.NoError(t, err)

			out, err := json.MarshalIndent(test.element, "", "  ")
			require.NoError(t, err)

			_, err = file.Write(out)
			require.NoError(t, err)
		})
	}
}

func TestGenRuntime(t *testing.T) {
	element := &runtime.Configuration{}

	// generator.Generate(element)

	hydration := NewHydration()
	hydration.Fill(element)

	cleanMiddlewaresRuntime(element)
	// cleanServers(element)

	file, err := os.Create("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/api_runtimeConfiguration.json")
	require.NoError(t, err)

	out, err := json.MarshalIndent(element, "", "  ")
	require.NoError(t, err)

	_, err = file.Write(out)
	require.NoError(t, err)
}

// func TestGenMiddlewares(t *testing.T) {
// 	file, err := os.Open("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/api_runtimeConfiguration.json")
// 	require.NoError(t, err)
//
// 	elementRuntime := &runtime.Configuration{}
// 	err = json.NewDecoder(file).Decode(elementRuntime)
// 	require.NoError(t, err)
//
// 	results := make([]api.MiddlewareRepresentation, 0, len(elementRuntime.Middlewares))
//
// 	for name, mi := range elementRuntime.Middlewares {
// 		results = append(results, api.NewMiddlewareRepresentation(name, mi))
// 	}
//
// 	output, err := os.Create("/home/ldez/sources/go/src/github.com/containous/traefik/docs/content/reference/dynamic-configuration/api_middlewares.json")
// 	require.NoError(t, err)
//
// 	encoder := json.NewEncoder(output)
// 	encoder.SetIndent("", "  ")
// 	err = encoder.Encode(results)
// 	require.NoError(t, err)
// }
