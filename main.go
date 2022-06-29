package main

import (
	"context"
	"flag"
	funcsdk "github.com/numaproj/numaflow-go/function"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const defaultURL = "https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&fields=regularMarketPrice,currency&symbols="
const defaultTicker = ""
const defaultMaxWaitTime = 3

//zap
func (q *queryHandler) handle(ctx context.Context, key, msg []byte) (funcsdk.Messages, error) {
	baseUrl := q.url
	requestMaxWaitTime := time.Duration(q.reqMaxWaitTime) * time.Second
	retMes := funcsdk.MessagesBuilder()
	ctx, cancel := context.WithTimeout(context.Background(), requestMaxWaitTime)
	defer cancel()
	//var messageReturn []byte
	for _, tick := range q.tickers {
		//Specific to current implementation
		url := baseUrl + tick
		data, err := q.processHttp(ctx, url)
		if err != nil {
			log.Print("Error Received %s", err)
		}
		log.Print("Received for tick: %s", tick)
		retMes = retMes.Append(funcsdk.MessageTo(tick, data))
	}
	return retMes, nil
}

func main() {
	url := flag.String("url", defaultURL, "URL to Query")
	tickers := flag.String("tickers", defaultTicker, "Comma seperated ticker names")
	flag.Parse()
	ticks := strings.Split(*tickers, ",")
	q := NewQueryHandler(1, defaultMaxWaitTime)
	q.url = *url
	q.tickers = ticks
	funcsdk.Start(context.Background(), q.handle)
}
