package libs

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Parallel struct {
	ProcessNum int
}

func NewParallel(config *Config) *Parallel {
	s := &Parallel{}
	s.ProcessNum = config.Parallel.ProcessNum
	return s
}

func (s Parallel) Run(urls *URLs) error {
	// split url
	splitUrls := urls.Split(s.ProcessNum)
	// create http client pool
	pool := NewHttpClientPool(s.ProcessNum, urls.Host)
	// 初期化
	result := NewResult(urls.Count)
	ctx := context.Background()
	pipe := make(chan int, 50)
	// ベンチマーク開始
	go parallelBenchmark(ctx, pipe, pool, splitUrls)
	// 結果を収集
	for i := 0; i < urls.Count; i++ {
		res := <-pipe
		if res == -1 {
			// ベンチマークタイムアウト
			break
		}
		result.AddCount(res)
	}
	// ベンチマーク終了
	result.Finish()
	// 正常終了
	return nil
}

func parallelBenchmark(ctx context.Context, wPipe chan<- int, pool *HttpClientPool, urls []*URLs) {
	childCtx, cancel := context.WithCancel(ctx)
	for i, httpClient := range pool.Clients {
		go sequentialBenchmark(childCtx, wPipe, httpClient, urls[i])
	}
	// 10秒でタイムアウト
	time.Sleep(10 * time.Second)
	cancel()
	wPipe <- -1
}

func sequentialBenchmark(ctx context.Context, wPipe chan<- int, httpClient http.Client, urls *URLs) {
	// エラーカウント
	errorCount := 0
	maxErrorCount := 5
	// ループ処理
	for _, requestInfo := range urls.Data {
		select {
		case <-ctx.Done():
			fmt.Println("親プロセスが処理をキャンセルしました")
			return
		default:
			// エラーの数が maxErrorCount を上回った場合に処理を停止
			if errorCount > maxErrorCount {
				fmt.Println("エラーが多すぎるため処理を中止しました")
				wPipe <- -1
				return
			}
			// リクエスト
			resp, err := httpClient.Get(requestInfo.String())
			if err != nil {
				errorCount += 1
				wPipe <- 500
				continue
			}
			// keepalive 用にデータを読み込む
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				errorCount += 1
				wPipe <- 500
				continue
			}
			err = resp.Body.Close()
			// レスポンス集計
			wPipe <- resp.StatusCode
		}
	}
}
