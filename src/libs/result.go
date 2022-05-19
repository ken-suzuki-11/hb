package libs

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"time"
)

type Result struct {
	TotalCount    int
	OkCount       int
	NotFoundCount int
	ErrorCount    int
	OtherCount    int
	StartTime     time.Time
	Bar           *pb.ProgressBar
}

func NewResult(totalCount int) *Result {
	return &Result{
		TotalCount:    totalCount,
		OkCount:       0,
		OtherCount:    0,
		NotFoundCount: 0,
		ErrorCount:    0,
		StartTime:     time.Now(),
		Bar:           pb.StartNew(totalCount),
	}
}

func (r *Result) AddCount(statusCode int) {
	r.Bar.Increment()

	switch statusCode {
	case 200:
		r.OkCount += 1
	case 404:
		r.NotFoundCount += 1
	case 500:
		r.ErrorCount += 1
	default:
		r.OtherCount += 1
	}
}

func (r *Result) Finish() {
	end := time.Now()
	totalTime := end.Sub(r.StartTime).Seconds()

	r.Bar.FinishPrint("\n")

	fmt.Println("### ベンチマーク結果 ###")
	fmt.Printf("実行回数: %d回\n", r.TotalCount)
	fmt.Printf("実行時間: %4.2v秒\n", totalTime)
	fmt.Printf("平均処理速度: %2.2v秒\n", totalTime/float64(r.TotalCount))
	fmt.Printf("Status 200 : %d\n", r.OkCount)
	fmt.Printf("Status 404 : %d\n", r.NotFoundCount)
	fmt.Printf("Status Error : %d\n", r.ErrorCount)
	fmt.Printf("Status Other : %d\n", r.OtherCount)
}
