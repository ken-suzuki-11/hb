package libs

import (
	"fmt"
	"github.com/pkg/errors"
)

type Sequential struct{}

func NewSequential(config *Config) *Sequential {
	s := &Sequential{}
	return s
}

func (s Sequential) Run(urls *URLs) error {
	// create http client pool
	pool := NewHttpClientPool(1, urls.Host)
	// Bench
	err := s.benchmark(pool, urls)
	if err != nil {
		fmt.Println("Error : ベンチマークに失敗しました")
		return errors.Wrap(err, "benchmark")
	}
	// 正常終了
	return nil
}

func (s Sequential) benchmark(pool *HttpClientPool, urls *URLs) error {
	// エラーカウント
	maxErrorCount := 5
	// 初期化
	result := NewResult(urls.Count)

	// ループ処理
	for i, requestInfo := range urls.Data {
		// エラーの数が maxErrorCount を上回った場合に処理を停止
		if result.ErrorCount > maxErrorCount {
			return errors.New("error count exceed")
		}
		// URLを生成
		targetUrl := fmt.Sprintf("%s://%s%s", requestInfo.Scheme, requestInfo.Host, requestInfo.Path)
		// pool内のhttpクライアントを利用してGETリクエスト
		index := i % pool.Num
		resp, err := pool.Clients[index].Get(targetUrl)
		if err != nil {
			result.ErrorCount += 1
			continue
		}
		// keepalive 用にデータを読み込む
		/*
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			result.ErrorCount += 1
			continue
		}
		err = resp.Body.Close()
		*/
		// ステータスコードのカウントをインクリメント
		result.AddCount(resp.StatusCode)
	}

	// ベンチマーク終了
	result.Finish()

	return nil
}
