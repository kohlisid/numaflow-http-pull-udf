package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"time"
)

func (q *queryHandler) processHttp(ctx context.Context, url string) ([]byte, error) {
	var data []byte
	time.Sleep(100 * time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Print("Cannot create request: %s\n", err)
		return nil, err
	}
	rsp, err := q.client.Do(req)
	if err != nil {
		log.Print("Error in req %s", err)
		return nil, err
	}
	if rsp != nil {
		defer rsp.Body.Close()
	}
	data, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Print("Error in reading req %s", err)
		return nil, err
	}
	return data, nil
}
