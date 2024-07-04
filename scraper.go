package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"slices"
	"strconv"

	"golang.org/x/net/html"
)

type Port struct {
	rxPkt  float64
	txPkt  float64
	crcPkt float64
}

func marshalPorts(readings []string) map[int]Port {
	var portReadings = make(map[int]Port)
	portNumber := 1
	var port Port
	for i := 0; i < len(readings); i++ {
		if i%3 == 0 {
			port = Port{}
			pr, _ := strconv.ParseInt(readings[i], 16, 0)
			port.rxPkt = float64(pr)
		} else if i%3 == 1 {
			pr, _ := strconv.ParseInt(readings[i], 16, 0)
			port.txPkt = float64(pr)
		} else if i%3 == 2 {
			pr, _ := strconv.ParseInt(readings[i], 16, 0)
			port.crcPkt = float64(pr)
			portReadings[portNumber] = port
			portNumber++
		}
	}
	return portReadings
}

func extractReadings(body io.Reader) ([]string, error) {
	var buf bytes.Buffer
	bod := io.TeeReader(body, &buf)
	doc, err := html.Parse(bod)
	if err != nil {
		return nil, err
	}
	bod2, _ := io.ReadAll(&buf)
	if bytes.Contains(bod2, []byte("RedirectToLoginPage")) {
		return nil, errors.New("redirect to loginpage")
	}
	var readings []string
	var page func(*html.Node)
	page = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, a := range n.Attr {
				if a.Key == "value" {
					readings = append(readings, a.Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			page(c)
		}
	}
	page(doc)
	if len(readings) == 0 {
		return nil, errors.New("empty extract page")
	}
	// skip the first two entries
	readings = slices.Delete(readings, 0, 2)

	return readings, nil
}

func (sc *scrapeClient) scrapeMetrics() (map[int]Port, error) {
	rand, err := sc.getRand()
	if err != nil {
		return nil, err
	}
	err = sc.login(rand)
	if err != nil {
		return nil, err
	}
	target, err := sc.urlBuilder("portStats.htm")
	if err != nil {
		return nil, err
	}
	statisticsRequest, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}

	statisticsResponse, err := sc.client.Do(statisticsRequest)
	if err != nil {
		return nil, err
	}
	defer statisticsResponse.Body.Close()

	readings, err := extractReadings(statisticsResponse.Body)
	if err != nil {
		return nil, err
	}

	return marshalPorts(readings), nil
}

func (sc *scrapeClient) crunchMetrics() error {
	nm, err := sc.scrapeMetrics()
	if err != nil {
		return err
	}

	sc.metric = nm

	for port := range sc.metric {
		rxPktMetric.WithLabelValues(strconv.Itoa(port), sc.location).Set(sc.metric[port].rxPkt)
		txPktMetric.WithLabelValues(strconv.Itoa(port), sc.location).Set(sc.metric[port].txPkt)
		crcPktMetric.WithLabelValues(strconv.Itoa(port), sc.location).Set(sc.metric[port].crcPkt)
	}

	return nil
}
