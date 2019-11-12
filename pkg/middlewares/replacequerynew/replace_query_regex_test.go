package replacequerynew

import (
	"context"
	"net/http"
	"testing"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceQueryRegex_onQuery_onQuery(t *testing.T) {
	testCases := []struct {
		desc               string
		request            string
		config             dynamic.ReplaceQueryNewRegex
		expectedQuery      string
		expectedRequestURI string
	}{
		{
			desc:    "no query no match",
			request: "/foo",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `(.+)`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "",
			expectedRequestURI: "/foo",
		},
		{
			desc:    "no query but match",
			request: "/foo",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `(.*)`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/foo?bar=baz",
		},
		{
			desc:    "no query but match, no path",
			request: "",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `(.*)`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/?bar=baz",
		},
		{
			desc:    "remove query parameter",
			request: "/foo?remove=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `.*(.*)`, // greedy leaves nothing
				ReplaceOn:   onQuery,
				Replacement: "$1",
			},
			expectedQuery:      "",
			expectedRequestURI: "/foo",
		},
		{
			desc:    "overwrite query parameters",
			request: "/foo?dropped=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `.*`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/foo?bar=baz",
		},
		{
			desc:    "append query parameter",
			request: "/foo?keep=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `(.*)`,
				ReplaceOn:   onQuery,
				Replacement: "$1&bar=baz",
			},
			expectedQuery:      "keep=yes&bar=baz",
			expectedRequestURI: "/foo?keep=yes&bar=baz",
		},
		{
			desc:    "modify query parameter",
			request: "/foo?@=a",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onQuery,
				Regex:       `@=a`,
				ReplaceOn:   onQuery,
				Replacement: "a=A",
			},
			expectedQuery:      "a=A",
			expectedRequestURI: "/foo?a=A",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var actualQuery, requestURI string
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestURI = r.RequestURI
				actualQuery = r.URL.RawQuery
			})

			handler, err := New(context.Background(), next, test.config, "foo-replace-query-regexp")
			require.NoError(t, err)

			req := testhelpers.MustNewRequest(http.MethodGet, "http://localhost"+test.request, nil)
			req.RequestURI = test.request

			handler.ServeHTTP(nil, req)

			assert.Equal(t, test.expectedQuery, actualQuery)
			assert.Equal(t, test.expectedRequestURI, requestURI)
		})
	}
}

func TestReplaceQueryRegex_onPathAndQuery_onQuery(t *testing.T) {
	testCases := []struct {
		desc               string
		request            string
		config             dynamic.ReplaceQueryNewRegex
		expectedQuery      string
		expectedRequestURI string
	}{
		{
			desc:    "add query parameter",
			request: "/foo",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/foo?bar=baz",
		},
		{
			desc:    "add query parameter, no path",
			request: "",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/?bar=baz",
		},
		{
			desc:    "remove query parameter",
			request: "/foo?remove=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*(.*)`, // greedy leaves nothing
				ReplaceOn:   onQuery,
				Replacement: "$1",
			},
			expectedQuery:      "",
			expectedRequestURI: "/foo",
		},
		{
			desc:    "overwrite query parameters",
			request: "/foo?dropped=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*`,
				ReplaceOn:   onQuery,
				Replacement: "bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/foo?bar=baz",
		},
		{
			desc:    "append query parameter",
			request: "/foo?keep=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^/foo\?(.*)$`,
				ReplaceOn:   onQuery,
				Replacement: "$1&bar=baz",
			},
			expectedQuery:      "keep=yes&bar=baz",
			expectedRequestURI: "/foo?keep=yes&bar=baz",
		},
		{
			desc:    "modify query parameter",
			request: "/foo?a=a",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^/foo\?a=a$`,
				ReplaceOn:   onQuery,
				Replacement: "a=A",
			},
			expectedQuery:      "a=A",
			expectedRequestURI: "/foo?a=A",
		},
		{
			desc:    "use path component as new query parameters",
			request: "/foo/animal/CAT/food/FISH?keep=no",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^/foo/animal/([^/]+)/food/([^?]+)(\?.*)?$`,
				ReplaceOn:   onQuery,
				Replacement: "animal=$1&food=$2",
			},
			expectedQuery:      "animal=CAT&food=FISH",
			expectedRequestURI: "/foo/animal/CAT/food/FISH?animal=CAT&food=FISH",
		},
		{
			desc:    "use path component as new query parameters, keep existing query params",
			request: "/foo/animal/CAT/food/FISH?keep=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^/foo/animal/([^/]+)/food/([^/]+)\?(.*)$`,
				ReplaceOn:   onQuery,
				Replacement: "$3&animal=$1&food=$2",
			},
			expectedQuery:      "keep=yes&animal=CAT&food=FISH",
			expectedRequestURI: "/foo/animal/CAT/food/FISH?keep=yes&animal=CAT&food=FISH",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var actualQuery, requestURI string
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestURI = r.RequestURI
				actualQuery = r.URL.RawQuery
			})

			handler, err := New(context.Background(), next, test.config, "foo-replace-query-regexp")
			require.NoError(t, err)

			req := testhelpers.MustNewRequest(http.MethodGet, "http://localhost"+test.request, nil)
			req.RequestURI = test.request

			handler.ServeHTTP(nil, req)

			assert.Equal(t, test.expectedQuery, actualQuery)
			assert.Equal(t, test.expectedRequestURI, requestURI)
		})
	}
}

func TestReplaceQueryRegex_onPathAndQuery_onPathAndQuery(t *testing.T) {
	testCases := []struct {
		desc               string
		request            string
		config             dynamic.ReplaceQueryNewRegex
		expectedQuery      string
		expectedRequestURI string
	}{
		{
			desc:    "add query parameter",
			request: "/foo",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "/foo?bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/foo?bar=baz",
		},
		{
			desc:    "add query parameter, no path",
			request: "",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "?bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/?bar=baz",
		},
		{
			desc:    "remove query parameter",
			request: "/foo?remove=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `.*(.*)`, // greedy leaves nothing
				ReplaceOn:   onPathAndQuery,
				Replacement: "/foo",
			},
			expectedQuery:      "",
			expectedRequestURI: "/foo",
		},
		{
			desc:    "overwrite query parameters",
			request: "/foo?dropped=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `(/foo)\?.*`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "$1?bar=baz",
			},
			expectedQuery:      "bar=baz",
			expectedRequestURI: "/foo?bar=baz",
		},
		{
			desc:    "append query parameter",
			request: "/foo?keep=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^(/foo)\?(.*)$`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "$1?$2&bar=baz",
			},
			expectedQuery:      "keep=yes&bar=baz",
			expectedRequestURI: "/foo?keep=yes&bar=baz",
		},
		{
			desc:    "modify query parameter",
			request: "/foo?a=a",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^(/foo)\?a=a$`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "$1?a=A",
			},
			expectedQuery:      "a=A",
			expectedRequestURI: "/foo?a=A",
		},
		{
			desc:    "use path component as new query parameters",
			request: "/foo/animal/CAT/food/FISH?keep=no",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^(/foo)/animal/([^/]+)/food/([^?]+)(\?.*)?$`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "$1?animal=$2&food=$3",
			},
			expectedQuery:      "animal=CAT&food=FISH",
			expectedRequestURI: "/foo?animal=CAT&food=FISH",
		},
		{
			desc:    "use path component as new query parameters, keep existing query params",
			request: "/foo/animal/CAT/food/FISH?keep=yes",
			config: dynamic.ReplaceQueryNewRegex{
				MatchOn:     onPathAndQuery,
				Regex:       `^(/foo)/animal/([^/]+)/food/([^/]+)\?(.*)$`,
				ReplaceOn:   onPathAndQuery,
				Replacement: "$1?$4&animal=$2&food=$3",
			},
			expectedQuery:      "keep=yes&animal=CAT&food=FISH",
			expectedRequestURI: "/foo?keep=yes&animal=CAT&food=FISH",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var actualQuery, requestURI string
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestURI = r.RequestURI
				actualQuery = r.URL.RawQuery
			})

			handler, err := New(context.Background(), next, test.config, "foo-replace-query-regexp")
			require.NoError(t, err)

			req := testhelpers.MustNewRequest(http.MethodGet, "http://localhost"+test.request, nil)
			req.RequestURI = test.request

			handler.ServeHTTP(nil, req)

			assert.Equal(t, test.expectedQuery, actualQuery)
			assert.Equal(t, test.expectedRequestURI, requestURI)
		})
	}
}

func TestReplaceQueryRegex_error(t *testing.T) {
	testCases := []struct {
		desc          string
		request       string
		config        dynamic.ReplaceQueryNewRegex
		expectedQuery string
	}{
		{
			desc:    "bad regex",
			request: "/foo",
			config: dynamic.ReplaceQueryNewRegex{
				Regex:       `(?!`,
				Replacement: "",
			},
			expectedQuery: "",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			_, err := New(context.Background(), next, test.config, "foo-replace-query-regexp")
			require.Error(t, err)
		})
	}
}
