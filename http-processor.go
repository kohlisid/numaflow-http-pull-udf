package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
)

func (q *queryHandler) processHttp(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Panicf("Cannot create request: %s\n", err)
	}
	rsp, err := q.client.Do(req)
	if err != nil {
		log.Println("Issues: %s\n", err)
		return nil, err
	}
	if rsp != nil {
		defer rsp.Body.Close()
	}
	data, err := ioutil.ReadAll(rsp.Body)
	return data, err

}
