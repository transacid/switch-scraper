package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type scrapeClient struct {
	client   *http.Client
	password string
	remote   string
	location string
	logger   *slog.Logger
	metric   map[int]Port
}

func NewScrapeClient() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}

	return client, nil
}

func (sc *scrapeClient) urlBuilder(path string) (string, error) {
	url, err := url.Parse(fmt.Sprintf("http://%s/%s", sc.remote, path))
	if err != nil {
		return "nil", err
	}
	return url.String(), nil
}

func (sc *scrapeClient) getRand() (string, error) {
	target, err := sc.urlBuilder("login.htm")
	if err != nil {
		return "", err
	}
	loginPageRequest, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return "", err
	}
	loginPageResponse, err := sc.client.Do(loginPageRequest)
	if err != nil {
		return "", err
	}
	defer loginPageResponse.Body.Close()

	rand, err := extractRand(loginPageResponse.Body)
	if err != nil {
		return "", err
	}
	return rand, nil
}

func (sc *scrapeClient) login(rand string) error {
	passwordHash := hashPassword(sc.password, rand)
	target, err := sc.urlBuilder("login.cgi")
	if err != nil {
		return err
	}
	form := url.Values{}
	form.Add("password", passwordHash)
	loginResponse, err := sc.client.PostForm(target, form)

	if err != nil {
		return err
	}

	defer loginResponse.Body.Close()

	_, err = io.ReadAll(loginResponse.Body)
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	return nil
}

func extractRand(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var readings []html.Attribute
	var page func(*html.Node)
	page = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "rand" {
					readings = append(readings, n.Attr...)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			page(c)
		}
	}
	page(doc)

	var rand string
	for _, element := range readings {
		if element.Key == "value" {
			rand = element.Val
		}
	}

	return rand, nil
}
