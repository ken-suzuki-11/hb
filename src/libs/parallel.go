package libs

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BenchmarkTimeout int = -1
	BenchMarkError   int = -2
)

type ParallelBenchmark struct {
	ProcessNum int
	TotalCount int
	Timeout    int
	HttpPool []http.Client
	SplitURLs  []*URLs
	Result *Result
}

func NewParallelBenchmark(config *Config, urls *URLs) *ParallelBenchmark {
	p := &ParallelBenchmark{}
	// 並列数初期化
	p.ProcessNum = config.Parallel.ProcessNum
	// 実行回数初期化
	p.TotalCount = urls.Count
	// タイムアウト値初期化
	p.Timeout = config.Common.Timeout
	// HTTP Client Pool 初期化
	tool := NewHttpClientPoolTool(config)
	p.HttpPool = tool.CreatePool(config.Parallel.ProcessNum, urls.Host)
	// URLの分割
	p.SplitURLs = urls.Split(config.Parallel.ProcessNum)
	// 結果の初期化
	p.Result = NewResult(urls.Count)
	return p
}

func (p ParallelBenchmark) Run() error {
	// コンテキスト作成
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Result Reset
	p.Result.StartTimeReset()
	// ベンチマーク開始
	rPipe := p.benchmark(childCtx)
	// 結果を収集
	for i := 0; i < p.TotalCount; i++ {
		res := <-rPipe
		if res == BenchmarkTimeout {
			fmt.Println("ベンチマークタイムアウト")
			cancel()
			break
		} else if res == BenchMarkError {
			fmt.Println("ベンチマークエラー")
			cancel()
			break
		} else {
			p.Result.AddCount(res)
		}
	}
	// ベンチマーク終了
	p.Result.Finish()
	// 正常終了
	return nil
}

func (p ParallelBenchmark) benchmark(ctx context.Context) <-chan int {
	// 結果通知用チャンネル
	pipe := make(chan int)
	// ベンチマーク開始
	go p.parallelProcess(ctx, pipe)
	return pipe
}

func (p ParallelBenchmark) parallelProcess(ctx context.Context, wPipe chan<- int) {
	// 並列処理
	for i, httpClient := range p.HttpPool {
		go p.sequentialProcess(ctx, wPipe, httpClient, p.SplitURLs[i])
	}
	// タイムアウト
	time.Sleep(time.Duration(p.Timeout) * time.Second)
	wPipe <- BenchmarkTimeout
}

func (p ParallelBenchmark) sequentialProcess(ctx context.Context, wPipe chan<- int, httpClient http.Client, urls *URLs) {
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
				wPipe <- BenchMarkError
				return
			}
			// リクエスト
			resp, err := httpClient.Get(requestInfo.String())
			if err != nil {
				errorCount += 1
				wPipe <- HttpError
				continue
			}
			// keepalive 用にデータを読み込む
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				errorCount += 1
				wPipe <- HttpError
				continue
			}
			err = resp.Body.Close()
			// レスポンス集計
			wPipe <- resp.StatusCode
		}
	}
}
