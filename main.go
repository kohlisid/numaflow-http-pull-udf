package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	funcsdk "github.com/numaproj/numaflow-go/function"
)

const defaultURL = "https://query1.finance.yahoo.com/v7/finance/quote?lang=en-US&region=US&corsDomain=finance.yahoo.com&fields=regularMarketPrice,currency&symbols="

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
		data, err := q.processHttp(ctx, url)
		if err != nil {
			fmt.Println("Error received %s", err)
		}
		messageReturn = append(messageReturn, data...)
		fmt.Println("DEBUG MESSAGE 5 : ", string(data), " and ", string(msg))
	}
	return funcsdk.MessagesBuilder().Append(funcsdk.MessageToAll(messageReturn)), nil
}

func main() {
	q := NewQueryHandler(1)
	url := flag.String("url", defaultURL, "URL to Query")
	tickers := flag.String("tickers", "", "Comma seperated ticker names")
	flag.Parse()
	ticks := strings.Split(*tickers, ",")
	q.url = *url
	q.tickers = ticks

	fmt.Println("DEBUG MESSAGE: ENTRY")
	funcsdk.Start(context.Background(), q.handle)
}
