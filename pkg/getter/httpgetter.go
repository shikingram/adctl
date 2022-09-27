package getter

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type HTTPGetter struct {
	opts options
}

func (g *HTTPGetter) Get(href string, options ...Option) (*bytes.Buffer, error) {
	for _, opt := range options {
		opt(&g.opts)
	}
	return g.get(href)
}

func (g *HTTPGetter) get(href string) (*bytes.Buffer, error) {
	req, err := http.NewRequest(http.MethodGet, href, nil)
	if err != nil {
		return nil, err
	}

	if g.opts.userAgent != "" {
		req.Header.Set("User-Agent", g.opts.userAgent)
	}

	u1, err := url.Parse(g.opts.url)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse getter URL")
	}
	u2, err := url.Parse(href)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse URL getting from")
	}

	if g.opts.passCredentialsAll || (u1.Scheme == u2.Scheme && u1.Host == u2.Host) {
		if g.opts.username != "" && g.opts.password != "" {
			req.SetBasicAuth(g.opts.username, g.opts.password)
		}
	}

	client, err := g.httpClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to fetch %s : %s", href, resp.Status)
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)
	return buf, err
}

func NewHTTPGetter(options ...Option) (Getter, error) {
	var client HTTPGetter

	for _, opt := range options {
		opt(&client.opts)
	}

	return &client, nil
}

func (g *HTTPGetter) httpClient() (*http.Client, error) {
	transport := &http.Transport{
		DisableCompression: true,
		Proxy:              http.ProxyFromEnvironment,
	}

	if g.opts.insecureSkipVerifyTLS {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		} else {
			transport.TLSClientConfig.InsecureSkipVerify = true
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   g.opts.timeout,
	}

	return client, nil
}
