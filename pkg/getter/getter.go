package getter

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
	"github.com/shikingram/adctl/pkg/cli"
)

type options struct {
	url                   string
	insecureSkipVerifyTLS bool
	username              string
	password              string
	passCredentialsAll    bool
	userAgent             string
	timeout               time.Duration
}

type Option func(*options)

func WithBasicAuth(username, password string) Option {
	return func(opts *options) {
		opts.username = username
		opts.password = password
	}
}

type Providers []Provider

func (p Providers) ByScheme(scheme string) (Getter, error) {
	for _, pp := range p {
		if pp.Provides(scheme) {
			return pp.New()
		}
	}
	return nil, errors.Errorf("scheme %q not supported", scheme)
}

type Provider struct {
	Schemes []string
	New     Constructor
}

func (p Provider) Provides(scheme string) bool {
	for _, i := range p.Schemes {
		if i == scheme {
			return true
		}
	}
	return false
}

type Constructor func(options ...Option) (Getter, error)

type Getter interface {
	Get(url string, options ...Option) (*bytes.Buffer, error)
}

func WithURL(url string) Option {
	return func(opts *options) {
		opts.url = url
	}
}

func All(settings *cli.EnvSettings) Providers {
	result := Providers{httpProvider}
	return result
}

var httpProvider = Provider{
	Schemes: []string{"http", "https"},
	New:     NewHTTPGetter,
}
