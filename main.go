package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	funcsdk "github.com/numaproj/numaflow-go/function"
)

type queryHandler struct {
	client *http.Client
}

func NewHTTPUDF() *queryHandler {

	//var httpClient *http.Client
	//// https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance/
	//httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	//// all our connects are loop back
	//httpTransport.MaxIdleConns = 100
	//httpTransport.MaxConnsPerHost = 100
	//httpTransport.MaxIdleConnsPerHost = 100
	//httpTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
	//	return net.Dial("unix", socketPath)
	//}
	//httpClient = &http.Client{
	//	Transport: httpTransport,
	//}
	//
	//return &queryHandler{
	//	client: httpClient,
	//}

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
	var url = os.Getenv("HTTP_URL")
	const RequestMaxWaitTime = 5 * time.Second
	fmt.Println("DEBUG MESSAGE: HANDLE ", url)
	ctx, cancel := context.WithTimeout(context.Background(), RequestMaxWaitTime)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Panicf("Cannot create request: %s\n", err)
	}
	fmt.Println("DEBUG MESSAGE: DO")
	rsp, err := q.client.Do(req)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	data, err := ioutil.ReadAll(rsp.Body)
	fmt.Println("DEBUG MESSAGE 5 : ", string(data), " and ", string(msg))
	return funcsdk.MessagesBuilder().Append(funcsdk.MessageToAll(data)), nil
}

func main() {
	q := NewHTTPUDF()
	fmt.Println("DEBUG MESSAGE: ENTRY")
	funcsdk.Start(context.Background(), q.handle)
}
