package main

import (
	"net"
	"net/http"
	"time"
)

type queryHandler struct {
	client         *http.Client
	url            string
	tickers        []string
	reqMaxWaitTime int
}

func NewQueryHandler(waitTime int, maxWaitTime int) *queryHandler {
	ConnectMaxWaitTime := time.Duration(waitTime) * time.Second

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}
	return &queryHandler{
		client:         &client,
		reqMaxWaitTime: maxWaitTime,
	}
}
