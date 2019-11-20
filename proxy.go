package main

import (
	"bytes"
	"context"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var (
	client *http.Client

)

func init() {
	dialer := &net.Dialer{
		Timeout:   connTimeout,
		KeepAlive: connKeepAlive,
	}
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, addr)
			},
			MaxIdleConnsPerHost: connMaxIdle,
			DisableCompression:  true,
		},
	}
}

func handle(c echo.Context) error {
	pc := c.(*PandaContext)
	req := c.Request()

	// parse request
	proxyReq, err := http.NewRequest(req.Method, "https://exhentai.org"+req.URL.RequestURI(), req.Body)
	if err != nil {
		return c.String(400, "bad request")
	}

	// fix header
	header := &req.Header
	proxyReq.Header = header.Clone()
	proxyReq.Header.Del("Authorization")
	proxyReq.Header.Del("Cookie")
	proxyReq.Header.Del("Accept-Encoding")
	proxyReq.Header.Del("Referer")
	if proxyReq.Header.Get("Origin") != "" {
		proxyReq.Header.Set("Origin", "https://exhentai.org")
	}
	for _, c := range req.Cookies() {
		if c.Name == "ipb_member_id" || c.Name == "ipb_pass_hash" {
			continue
		}
		proxyReq.AddCookie(c)
	}
	proxyReq.AddCookie(&http.Cookie{Name: "ipb_member_id", Value: pc.MemberID})
	proxyReq.AddCookie(&http.Cookie{Name: "ipb_pass_hash", Value: pc.PassHash})
	proxyReq.ContentLength = req.ContentLength

	// send request
	resp, err := client.Do(proxyReq) // todo: connection pool
	if err != nil {
		return c.String(500, "server error")
	}

	// set header
	r := c.Response()
	for k, _ := range resp.Header {
		r.Header().Set(k, resp.Header.Get(k))
	}

	// fix header
	r.Header().Del("Content-Length")
	cookieString := resp.Header.Get("Set-Cookie")
	if cookieString != "" {
		cookieString = strings.ReplaceAll(cookieString, "; domain=.exhentai.org", "")
		r.Header().Set("Set-Cookie", cookieString)
	}

	// check content type
	contentType := strings.Split(resp.Header.Get("Content-Type"), ";")[0]
	if contentType[:4] != "text" {
		return c.Stream(resp.StatusCode, contentType, resp.Body)
	}

	// replace content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(500, "server error")
	}
	_ = resp.Body.Close()
	body = bytes.ReplaceAll(body, []byte("https://exhentai.org"), []byte(rootPath))
	body = bytes.ReplaceAll(body, []byte("exhentai.org"), []byte(rootHost))

	return c.Stream(resp.StatusCode, contentType, bytes.NewReader(body))
}
