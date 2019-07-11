package gen

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/label"
	"github.com/stretchr/testify/require"
)

func TestGenLabels(t *testing.T) {
	element := &dynamic.Configuration{}
	// generator.Generate(element)

	hydration := NewHydration()
	hydration.Fill(element)

	cleanMiddlewares(element)
	cleanHTTPServices(element)
	cleanTCPServices(element)
	cleanUDPServices(element)
	cleanServers(element)

	// flats, err := flag.Encode(element)
	// require.NoError(t, err)

	labels, err := label.EncodeConfiguration(element)
	require.NoError(t, err)

	var keys []string

	for k := range labels {
		if strings.HasPrefix(strings.ToLower(k), "traefik.tls.") {
			continue
		}
		keys = append(keys, k)
	}

	sort.Strings(keys)

	dockerLabels, err := os.Create(dynDockerLabels)
	require.NoError(t, err)

	marathonLabels, err := os.Create(dynMarathonLabels)
	require.NoError(t, err)

	for i, k := range keys {
		if labels[k] != "" {
			fmt.Fprintln(dockerLabels, `- "`+strings.ToLower(k)+`=`+labels[k]+`"`)

			if i == len(keys)-1 {
				fmt.Fprintln(marathonLabels, `"`+strings.ToLower(k)+`": "`+labels[k]+`"`)
			} else {
				fmt.Fprintln(marathonLabels, `"`+strings.ToLower(k)+`": "`+labels[k]+`",`)
			}
		}
	}

	// for _, v := range flats {
	// 	if strings.HasPrefix(strings.ToLower(v.Name), "tls.") {
	// 		continue
	// 	}
	// 	fmt.Fprintln(dockerLabels, `- "traefik.`+strings.ToLower(v.Name)+`=`+v.Default+`"`)
	// }

}
