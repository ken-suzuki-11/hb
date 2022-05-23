package libs

import (
	"crypto/tls"
	"net/http"
	"time"
)

type HttpClientPoolTool struct {
	DisableKeepAlive bool
}

func NewHttpClientPoolTool(config *Config) *HttpClientPoolTool {
	return &HttpClientPoolTool{
		DisableKeepAlive: config.Http.DisableKeepalive,
	}
}

func (h HttpClientPoolTool) CreatePool(num int, host string) []http.Client {
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
				DisableKeepAlives: h.DisableKeepAlive,
			},
		}
		pool = append(pool, httpClient)
	}
	return pool
}
