package scrape_stats

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
)

var ErrRequestForbidden = errors.New("request forbidden")

func requestGet(requestUrl string) (*http.Response, error) {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100)
	time.Sleep(time.Duration(r) * time.Millisecond)

	torProxy := "socks5://127.0.0.1:9050"
	torProxyUrl, err := url.Parse(torProxy)
	if err != nil {
		log.Fatal("Error parsing Tor proxy URL:", torProxy, ".", err)
	}
	torTransport := &http.Transport{Proxy: http.ProxyURL(torProxyUrl)}

	client := &http.Client{Transport: torTransport, Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", browser.Random())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Origin", "https://www.sofascore.com")
	req.Header.Set("Referer", "https://www.sofascore.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("TE", "trailers")

	return client.Do(req)
}
