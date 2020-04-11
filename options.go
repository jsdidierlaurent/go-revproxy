package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var DefaultParams = &Options{ProxyPort: 3000}

type (
	// Options contains command line options
	Options struct {
		// RemoteURL is the proxy target
		RemoteURL *url.URL

		// ProxyPort is the proxy port
		ProxyPort int

		// ProxyPrefix is add to the proxy adress
		// Ex: http://localhost:3000/prefix/
		ProxyPrefix string
	}
)

func NewOptions(remoteURL string, proxyPort int, proxyPrefix string) (*Options, error) {
	remoteURL = strings.TrimSuffix(remoteURL, "/")
	URL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, err
	}

	if proxyPort == 0 {
		return nil, errors.New("port must be > 0")
	}

	if proxyPrefix != "" {
		proxyPrefix = strings.Trim(proxyPrefix, "/")
		proxyPrefix = fmt.Sprintf("/%s", proxyPrefix)
	}

	return &Options{
		RemoteURL:   URL,
		ProxyPort:   proxyPort,
		ProxyPrefix: proxyPrefix,
	}, nil
}
