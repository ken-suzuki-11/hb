package main

import (
	"flag"
	"fmt"
	"hb/src/libs"
	"os"
)

func main() {
	var (
		configFile  string
		urlListFile string
	)
	flag.StringVar(&configFile, "c", "nil", "config file path")
	flag.StringVar(&urlListFile, "u", "nil", "url list file path")
	flag.Parse()
	// フラグの値をチェック
	if &configFile == nil || &urlListFile == nil {
		fmt.Println("Usage: ./prog -c config_path -u url_list_path")
		os.Exit(-1)
	}
	// 設定ファイルチェック
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		fmt.Println("Config File Not Found")
		os.Exit(-1)
	}
	// 設定読み込み
	config, err := libs.NewConfig(configFile)
	if err != nil {
		fmt.Println("Error : 設定の読み込みに失敗しました")
		fmt.Println(err)
		os.Exit(-1)
	}
	// デバッグの有無
	isDebug := config.Common.Debug
	// URLリストの読み込み
	urlTool := libs.NewURLsTool(urlListFile, config.Common.ListSizeLimit)
	urls, err := urlTool.Load()
	if err != nil {
		fmt.Println("Error : URLリストの読み込みに失敗しました")
		if isDebug {
			fmt.Printf("\n[StackTrace]\n%+v\n", err)
		}
		os.Exit(-1)
	}
	// ベンチマーク
	switch config.Function.Name {

	case "parallel":
		fmt.Println("parallel ベンチマーク")
		function := libs.NewParallelBenchmark(config,urls)
		err := function.Run()
		if err != nil {
			fmt.Println("エラーが発生しました")
			if isDebug {
				fmt.Printf("\n[StackTrace]\n%+v\n", err)
			}
		}

	default:
		fmt.Println("存在しない機能です")
	}
}
