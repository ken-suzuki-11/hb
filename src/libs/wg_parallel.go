package libs

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sync"
)

type WgParallel struct{
	ProcessNum int
}

func NewWgParallel(config *Config) *WgParallel {
	s := &WgParallel{}
	s.ProcessNum = config.WgParallel.ProcessNum
	return s
}

func (s WgParallel) Run(urls *URLs) error {
	// split url
	splitUrls := urls.Split(s.ProcessNum)
	// create http client pool
	pool := NewHttpClientPool(s.ProcessNum, urls.Host)
	// Bench
	
	err := s.benchmark(pool, urls.Count, splitUrls)
	if err != nil {
		fmt.Println("Error : ベンチマークに失敗しました")
		return errors.Wrap(err, "benchmark")
	}
	// 正常終了
	return nil
}

func (s WgParallel) benchmark(pool *HttpClientPool, totalCount int, urls []*URLs) error {

	// コンテキスト
	ctx := context.Background()

	// 初期化
	result := NewResult(totalCount)

	// 100バッファintチャンネル
	pipe := make(chan int, 100)

	wg := new(sync.WaitGroup)

	for i, httpClient := range pool.Clients {
		wg.Add(1)
		go s.sequentialLoop(wg, httpClient, urls[i], pipe)
	}

	for i:=0;i<totalCount;i++{
		res := <-pipe
		fmt.Println(res)
	}

	wg.Wait()

	// ベンチマーク終了
	result.Finish()

	return nil
}

func (s WgParallel)sequentialLoop(wg *sync.WaitGroup, httpClient http.Client, urls *URLs, wPipe chan<-int){
	defer wg.Done()

	// エラーカウント
	errorCount := 0
	maxErrorCount := 5
	// ループ処理
	for _, requestInfo := range urls.Data {
		// エラーの数が maxErrorCount を上回った場合に処理を停止
		if errorCount > maxErrorCount {
			return
		}
		// URLを生成
		targetUrl := fmt.Sprintf("%s://%s%s", requestInfo.Scheme, requestInfo.Host, requestInfo.Path)
		// pool内のhttpクライアントを利用してGETリクエスト
		resp, err := httpClient.Get(targetUrl)
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

		wPipe <- resp.StatusCode

		// ステータスコードのカウントをインクリメント
		//result.AddCount(resp.StatusCode)
	}
}