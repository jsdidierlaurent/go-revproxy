package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var port int
var prefix string
var remoteURL *url.URL

type ProxyHandler struct {
	p *httputil.ReverseProxy
	prefix string
}
func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL, _ = url.Parse(strings.TrimPrefix(r.URL.Path, ph.prefix))
	log.Print(r.URL)
	ph.p.ServeHTTP(w, r)
}

func main() {
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	http.Handle(fmt.Sprintf("%s/", prefix), &ProxyHandler{proxy, prefix})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func init() {
	// REMOTE_URL
	if len(os.Args) == 1 {
		fmt.Println("Usage: revproxy [OPTION]... REMOTE_URL")
		os.Exit(1)
	}

	var err error
	remote := os.Args[len(os.Args)-1]
	remoteURL, err = url.Parse(remote)
	if err != nil || (!strings.HasPrefix(remote, "http://") && !strings.HasPrefix(remote, "https://")) {
		fmt.Println("REMOTE_URL must be a valid url starting with http(s)://")
		os.Exit(1)
	}

	// FLAGS
	portFlag := flag.Int("port", 3000, "Reverse proxy port. (Required)")
	prefixFlag := flag.String("prefix", "", "route prefix (ex: localhost:3000/myPrefix/")
	flag.Parse()

	if portFlag == nil || *portFlag == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	port = *portFlag

	prefix = *prefixFlag
	if prefix != "" {
		prefix = strings.TrimPrefix(prefix, "/")
		prefix = strings.TrimSuffix(prefix, "/")
		prefix = fmt.Sprintf("/%s", prefix)
	}
}