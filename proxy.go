package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func newReverseProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse("http://" + target)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	director := func(req *http.Request) {
		original := req.Host
		go logRequest(req, original, target)

		req.URL.Scheme = "http"
		req.URL.Host = url.Host
		req.URL.Path = singleJoiningSlash(url.Path, req.URL.Path)
		if !strings.HasSuffix(req.URL.Path, "/") && req.URL.RawPath != "" {
			req.URL.RawPath = singleJoiningSlash(url.Path, req.URL.RawPath)
		}

		req.Host = url.Host
	}

	proxy.Director = director
	proxy.Transport = defaultTransport
	return proxy
}

func singleJoiningSlash(a, b string) string {
	a = strings.TrimRight(a, "/")
	b = strings.TrimLeft(b, "/")
	if a == "" {
		return "/" + b
	}
	return a + "/" + b
}

var defaultTransport http.RoundTripper = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	Proxy:           http.ProxyFromEnvironment,
}

func logRequest(req *http.Request, original string, target string) {
	uri := req.URL.RequestURI()
	if uri == "" {
		uri = "/"
	}

	log.Printf("[Epic302] %s → %s \"http://%s%s\"",
		original,
		target,
		target,
		uri,
	)
}
