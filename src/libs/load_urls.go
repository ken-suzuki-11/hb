package libs

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"os"
)

const maxCount int = 2000000000

type URLs struct {
	file  string
	limit int64
}

func NewURLs(file string, limit int64) *URLs {
	return &URLs{
		file:  file,
		limit: limit,
	}
}

func (u URLs) Load() (int, []*url.URL, error) {
	count := 0
	var data []*url.URL
	// File open
	file, err := u.openFile()
	if err != nil {
		fmt.Println("Error : ファイルのオープンに失敗しました")
		return count, data, errors.Wrap(err, "openFile")
	}
	// パスリストの読み込み
	count, listData, err := u.loadData(file)
	if err != nil {
		fmt.Println("Error : リストの読み込みに失敗しました")
		return count, data, errors.Wrap(err, "loadData")
	}
	// 正常終了
	return count, listData, nil
}

func (u URLs) openFile() (*os.File, error) {
	info, err := os.Stat(u.file)
	if os.IsNotExist(err) {
		fmt.Printf("Error : リストファイルが存在しません\n")
		return nil, errors.New("path list not found")
	}
	if info.Size() > u.limit {
		fmt.Println("Error : 読み込み容量オーバー")
		fmt.Printf("List Size : %dByte (%d)\n", info.Size(), u.limit)
		return nil, errors.New("file size too large")
	}
	file, err := os.Open(u.file)
	if err != nil {
		fmt.Printf("Error : ファイルオープンエラー : %s\n", u.file)
		return nil, errors.New("file open error")
	}
	return file, nil
}

func (u URLs) loadData(file *os.File) (int, []*url.URL, error) {
	// レスポンス用データ
	var data []*url.URL
	// Scanner作成
	fileScanner := bufio.NewScanner(file)
	// ホストチェック用
	hostMap := map[string]int{}
	count := 0
	// リスト読み込み
	for fileScanner.Scan() {
		line := fileScanner.Text()
		urlInfo, err := u.parseUrl(line)
		if err != nil {
			fmt.Printf("Error : Parseに失敗しました : %s\n", line)
			return 0, nil, errors.Wrap(err, "parse url")
		}
		_, ok := hostMap[urlInfo.Host]
		if !ok {
			hostMap[urlInfo.Host] = 1
		}
		data = append(data, urlInfo)
		count += 1
		if count > maxCount {
			fmt.Printf("Error : 処理可能な行数を超えました : %d\n", count)
			return 0, nil, errors.New("count too large")
		}
	}
	// ホストチェック
	if len(hostMap) != 1 {
		fmt.Println("Error : 複数のホストが含まれています")
		return 0, nil, errors.New("input host error")
	}
	return count, data, nil
}

func (u URLs) parseUrl(line string) (*url.URL, error) {
	urlInfo, err := url.Parse(line)
	if err != nil {
		return nil, err
	}
	if urlInfo.Scheme == "" {
		return nil, errors.New("scheme not found")
	}
	if urlInfo.Host == "" {
		return nil, errors.New("host not found")
	}
	if urlInfo.Path == "" {
		return nil, errors.New("path not found")
	}
	return urlInfo, nil
}
