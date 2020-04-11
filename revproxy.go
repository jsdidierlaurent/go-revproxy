package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type (
	ProxyHandler struct {
		proxy   *httputil.ReverseProxy
		options *Options
	}
)

var startup = `
Reverse proxy is running at:
  %s
  %s

`

func Run(o *Options) {
	// Creating RemoteURL proxy
	proxy := httputil.NewSingleHostReverseProxy(o.RemoteURL)

	// Bind this proxy on route with or without prefix
	http.Handle(o.getRoutePattern(), &ProxyHandler{proxy, o})

	// Starting http server
	Colorer.Printf(startup,
		Colorer.Blue(fmt.Sprintf("http://localhost:%d%s", o.ProxyPort, o.getRoutePattern())),
		Colorer.Blue(fmt.Sprintf("http://%s:%d%s", getNetworkIP(), o.ProxyPort, o.getRoutePattern())),
	)
	err := http.ListenAndServe(fmt.Sprintf(":%d", o.ProxyPort), nil)
	if err != nil {
		panic(err)
	}
}

func (p *Options) getRoutePattern() string {
	return fmt.Sprintf("%s/", p.ProxyPrefix)
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL, _ = url.Parse(strings.TrimPrefix(r.URL.Path, ph.options.ProxyPrefix))
	Colorer.Printf("[%6s] %s\n", Colorer.Green(r.Method), Colorer.Blue(fmt.Sprintf("%s%s", ph.options.RemoteURL, r.RequestURI)))
	ph.proxy.ServeHTTP(w, r)
}

func getNetworkIP() string {
	ip := "0.0.0.0"

	conn, err := net.Dial("udp", "255.255.255.255:80")
	if err == nil {
		defer conn.Close()
		networkIP := conn.LocalAddr().(*net.UDPAddr).IP
		ip = networkIP.To4().String()
	}

	return ip
}
