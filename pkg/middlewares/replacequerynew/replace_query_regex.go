package replacequerynew

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"github.com/containous/traefik/v2/pkg/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	typeName = "ReplacePathQueryRegex"
)

const (
	onQuery        = "query"
	onPathAndQuery = "pathandquery"
)

// replaceQueryRegex is a middleware used to replace the path and query of a URL request with a regular expression
type replaceQueryRegex struct {
	name        string
	matchOn     string
	regexp      *regexp.Regexp
	replaceOn   string
	replacement string
	next        http.Handler
}

// New creates a new replace path and query regex middleware.
func New(ctx context.Context, next http.Handler, config dynamic.ReplaceQueryNewRegex, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")

	exp, err := regexp.Compile(strings.TrimSpace(config.Regex))
	if err != nil {
		return nil, fmt.Errorf("error compiling regular expression %s: %w", config.Regex, err)
	}

	matchOn, err := getSourceType(config.MatchOn)
	if err != nil {
		return nil, fmt.Errorf("matchOn configuration error: %w", err)
	}

	replaceOn, err := getSourceType(config.ReplaceOn)
	if err != nil {
		return nil, fmt.Errorf("replaceOn configuration error: %w", err)
	}

	return &replaceQueryRegex{
		matchOn:     matchOn,
		regexp:      exp,
		replaceOn:   replaceOn,
		replacement: strings.TrimSpace(config.Replacement),
		next:        next,
		name:        name,
	}, nil
}

func (r *replaceQueryRegex) GetTracingInformation() (string, ext.SpanKindEnum) {
	return r.name, tracing.SpanKindNoneEnum
}

func (r *replaceQueryRegex) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	src := r.getSource(req)

	if r.regexp == nil || !r.regexp.MatchString(src) {
		r.next.ServeHTTP(rw, req)
		return
	}

	switch r.replaceOn {
	case onQuery:
		req.URL.RawQuery = r.regexp.ReplaceAllString(src, r.replacement)
	case onPathAndQuery:
		uri, err := url.Parse(r.regexp.ReplaceAllString(src, r.replacement))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		req.URL = uri
	}

	req.RequestURI = req.URL.RequestURI()

	r.next.ServeHTTP(rw, req)
}

func (r *replaceQueryRegex) getSource(req *http.Request) string {
	switch r.matchOn {
	case onQuery:
		return req.URL.RawQuery
	case onPathAndQuery:
		return req.RequestURI
	default:
		return req.URL.RawQuery
	}
}

func getSourceType(value string) (string, error) {
	if value == "" {
		return onQuery, nil
	}

	if strings.EqualFold(value, onQuery) {
		return onQuery, nil
	}

	if strings.EqualFold(value, onPathAndQuery) {
		return onPathAndQuery, nil
	}

	return "", fmt.Errorf("unsupported type: %s", value)
}
