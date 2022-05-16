package libs

import (
	"crypto/tls"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"time"
)

type Sequential struct{}

func NewSequential(config *Config) *Sequential {
	s := &Sequential{}
	return s
}

func (s *Sequential) Run(count int, urls []*url.URL) error {
	// Bench
	err := s.benchmark(count, urls)
	if err != nil {
		fmt.Println("Error : ベンチマークに失敗しました")
		return errors.Wrap(err, "benchmark")
	}
	// 正常終了
	return nil
}

func (s Sequential) benchmark(count int, data []*url.URL) error {
	firstData := data[0]
	host := firstData.Host

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

	successCount := 0
	notFoundCount := 0
	errorCount := 0
	unknownCount := 0

	fmt.Println("Benchmark Start")

	// プログレスバー初期化
	bar := pb.StartNew(count)
	// ベンチマーク開始
	// 開始時刻を取得
	start := time.Now()
	for _, requestInfo := range data {
		targetUrl := fmt.Sprintf("%s://%s%s", requestInfo.Scheme, requestInfo.Host, requestInfo.Path)
		resp, err := httpClient.Get(targetUrl)
		if err != nil {
			errorCount += 1
			fmt.Printf("%+v\n", err)
			return errors.Wrap(err, "http client get")
		}
		switch resp.StatusCode {
		case 200:
			successCount += 1
		case 404:
			notFoundCount += 1
		default:
			unknownCount += 1
		}
		bar.Increment()
	}

	// ベンチマーク終了
	bar.FinishPrint("Benchmark End\n")
	end := time.Now()
	totalTime := end.Sub(start).Seconds()

	fmt.Println("### ベンチマーク結果 ###")
	fmt.Printf("実行回数 : %d回\n", count)
	fmt.Printf("実行時間: %v秒\n", totalTime)
	fmt.Printf("平均処理速度: %v秒\n", totalTime/float64(count))

	fmt.Println(successCount)
	fmt.Println(notFoundCount)
	fmt.Println(errorCount)
	fmt.Println(unknownCount)
	return nil
}
