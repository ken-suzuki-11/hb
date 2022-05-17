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
	Count int
	Host  string
	Data  []*url.URL
}

func (u *URLs) Load(filepath string, limit int64) error {
	// File open
	file, err := u.openFile(filepath, limit)
	if err != nil {
		fmt.Println("Error : ファイルのオープンに失敗しました")
		return errors.Wrap(err, "openFile")
	}
	// パスリストの読み込み
	count, host, listData, err := u.loadData(file)
	if err != nil {
		fmt.Println("Error : リストの読み込みに失敗しました")
		return errors.Wrap(err, "loadData")
	}
	// 値を設定
	u.Count = count
	u.Host = host
	u.Data = listData
	// 正常終了
	return nil
}

func (u URLs) openFile(filepath string, limit int64) (*os.File, error) {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("Error : ファイルが存在しません : %s\n", filepath)
		return nil, errors.New("path list not found")
	}
	if info.Size() > limit {
		fmt.Println("Error : 読み込み容量オーバー")
		fmt.Printf("List Size : %dByte (%d)\n", info.Size(), limit)
		return nil, errors.New("file size too large")
	}
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error : ファイルオープンエラー : %s\n", filepath)
		return nil, errors.New("file open error")
	}
	return file, nil
}

func (u URLs) loadData(file *os.File) (int, string, []*url.URL, error) {
	var host string
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
			return 0, "", nil, errors.Wrap(err, "parse url")
		}
		_, ok := hostMap[urlInfo.Host]
		if !ok {
			host = urlInfo.Host
			hostMap[urlInfo.Host] = 1
		}
		data = append(data, urlInfo)
		count += 1
		if count > maxCount {
			fmt.Printf("Error : 処理可能な行数を超えました : %d\n", count)
			return 0, "", nil, errors.New("count too large")
		}
	}
	// ホストチェック
	if len(hostMap) != 1 {
		fmt.Println("Error : 複数のホストが含まれています")
		return 0, "", nil, errors.New("input host error")
	}
	return count, host, data, nil
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
