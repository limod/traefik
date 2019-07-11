package gen

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/env"
	"github.com/containous/traefik/v2/pkg/config/generator"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/stretchr/testify/require"
)

func TestGenEnv(t *testing.T) {
	element := &static.Configuration{}

	generator.Generate(element)

	flats, _ := env.Encode(element)

	for _, v := range flats {
		fmt.Println("`" + strings.ReplaceAll(v.Name, "[0]", "[n]") + "`:  ")
		if v.Default == "" {
			fmt.Println(v.Description)
		} else {
			fmt.Println(v.Description + " (Default: ```" + v.Default + "```)")
		}
		fmt.Println()

		// `TRAEFIK_ACCESSLOG`:
		//
		// - Access log settings. (Default: `false` )
	}
}

func TestGenEnv1(t *testing.T) {
	element := &static.Configuration{}

	generator.Generate(element)

	flats, err := env.Encode(element)
	require.NoError(t, err)

	file, err := os.Create(staticEnvRef)
	require.NoError(t, err)

	defer file.Close()

	for i, flat := range flats {
		// if flat.Hidden {
		// 	continue
		// }

		fmt.Fprintln(file, "`"+strings.ReplaceAll(flat.Name, "[0]", "[n]")+"`:  ")
		if flat.Default == "" {
			fmt.Fprintln(file, flat.Description)
		} else {
			fmt.Fprintln(file, flat.Description+" (Default: ```"+flat.Default+"```)")
		}

		if i < len(flats)-1 {
			fmt.Fprintln(file)
		}
	}
}
