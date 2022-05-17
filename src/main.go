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
	switch config.Common.Function {

	case "sequential":
		fmt.Println("シーケンシャルベンチマーク")
		function := libs.NewSequential(config)
		err := function.Run(urls)
		if err != nil {
			fmt.Println("シーケンシャルベンチマークでエラーが発生しました")
			if isDebug {
				fmt.Printf("\n[StackTrace]\n%+v\n", err)
			}
		}

	case "parallel":
		fmt.Println("WaitGroupパラレルベンチマーク")
		function := libs.NewParallel(config)
		err := function.Run(urls)
		if err != nil {
			fmt.Println("WaitGroupパラレルベンチマークでエラーが発生しました")
			if isDebug {
				fmt.Printf("\n[StackTrace]\n%+v\n", err)
			}
		}

	default:
		fmt.Println("存在しない機能です")
	}
}
