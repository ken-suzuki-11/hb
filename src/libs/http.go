package libs

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HttpClientPool struct {
	Num     int
	Clients []http.Client
}

func NewHttpClientPool(num int, host string) *HttpClientPool {
	client := HttpClientPool{}

	var pool []http.Client
	for i := 0; i < num; i++ {
		httpClient := http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
				TLSClientConfig: &tls.Config{
					ServerName: host,
				},
			},
		}
		pool = append(pool, httpClient)
	}
	client.Num = num
	client.Clients = pool

	return &client
}
