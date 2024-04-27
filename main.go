package main

import (
	"compress/gzip"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gofiber/fiber/v2"
)

var client *http.Client

func main() {

	proxy := os.Getenv("PROXY")
	var httpProxy func(*http.Request) (*url.URL, error)

	if proxy != "" {
		proxyUrl, _ := url.Parse(proxy)
		httpProxy = http.ProxyURL(proxyUrl)
	}

	client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 20 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       30 * time.Second,
			ForceAttemptHTTP2:     true,
			MaxConnsPerHost:       0,
			MaxIdleConnsPerHost:   10,
			MaxIdleConns:          0,
			Proxy:                 httpProxy,
		},
	}

	app := fiber.New(
		fiber.Config{
			Prefork: true,
		},
	)

	// Route for /votes?videoId=:videoId
	app.Get("/votes", handleQuery)

	// Route for /votes/:videoId
	app.Get("/votes/:videoId", handleParam)

	log.Fatal(app.Listen(":3000"))
}

func handleQuery(c *fiber.Ctx) error {
	videoId := c.Query("videoId")
	return getVotes(c, videoId)
}

func handleParam(c *fiber.Ctx) error {
	videoId := c.Params("videoId")
	return getVotes(c, videoId)
}

func getVotes(c *fiber.Ctx, videoId string) error {
	match, _ := regexp.Match("^([a-zA-Z0-9_-]{11})", []byte(videoId))

	if !match {
		return c.Status(400).SendString("Invalid video id")
	}

	for {
		req, _ := http.NewRequest("GET", "https://returnyoutubedislikeapi.com/Votes?videoId="+videoId+"&likeCount=", nil)

		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Accept-Language", "en-US,en;q=0.5")
		req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("Sec-Fetch-Dest", "empty")
		req.Header.Add("Sec-Fetch-Mode", "cors")
		req.Header.Add("Sec-Fetch-Site", "same-origin")

		resp, err := client.Do(req)

		if err != nil || resp.StatusCode == 429 {
			continue
		}

		defer resp.Body.Close()

		ce := resp.Header.Get("Content-Encoding")

		var stream io.Reader

		if ce == "gzip" {
			stream, _ = gzip.NewReader(resp.Body)
		} else if ce == "br" {
			stream = brotli.NewReader(resp.Body)
		} else {
			stream = resp.Body
		}

		return c.Status(resp.StatusCode).SendStream(stream)
	}
}
