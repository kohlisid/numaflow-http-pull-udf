package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	funcsdk "github.com/numaproj/numaflow-go/function"
)

type queryHandler struct {
	client  *http.Client
	url     string
	tickers []string
}

const defaultURL = "https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&fields=regularMarketPrice,currency&symbols="

func NewQueryHandler() *queryHandler {
	const ConnectMaxWaitTime = 1 * time.Second

	client := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: ConnectMaxWaitTime,
			}).DialContext,
		},
	}
	return &queryHandler{
		client: &client,
	}
}

//logres zap
func (q *queryHandler) handle(ctx context.Context, key, msg []byte) (funcsdk.Messages, error) {
	//var url = os.Getenv("HTTP_URL")
	baseUrl := q.url
	const RequestMaxWaitTime = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), RequestMaxWaitTime)
	defer cancel()
	var messageReturn []byte
	for _, tick := range q.tickers {
		//Specific to current implementaion
		url := baseUrl + tick
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		fmt.Println("DEBUG MESSAGE: HANDLE ", url)
		if err != nil {
			log.Panicf("Cannot create request: %s\n", err)
		}
		fmt.Println("DEBUG MESSAGE: DO")
		rsp, err := q.client.Do(req)
		if rsp != nil {
			defer rsp.Body.Close()
		}
		data, err := ioutil.ReadAll(rsp.Body)
		messageReturn = append(messageReturn, data...)
		fmt.Println("DEBUG MESSAGE 5 : ", string(data), " and ", string(msg))
	}
	return funcsdk.MessagesBuilder().Append(funcsdk.MessageToAll(messageReturn)), nil
}

func main() {
	q := NewQueryHandler()
	url := flag.String("url", defaultURL, "URL to Query")
	tickers := flag.String("tickers", "INTU", "Comma seperated ticker names")
	flag.Parse()
	ticks := strings.Split(*tickers, ",")
	q.url = *url
	q.tickers = ticks

	fmt.Println("DEBUG MESSAGE: ENTRY")
	funcsdk.Start(context.Background(), q.handle)
}
