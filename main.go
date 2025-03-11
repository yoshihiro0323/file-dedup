package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// ファイルグループは保持するファイルと削除候補を含む構造体
type FileGroup struct {
	keep     string   // 保持するファイル
	removals []string // 削除候補のファイル
}

// ソートアルゴリズムを適用してFileGroupを生成する
func createFileGroup(files []string) FileGroup {
	// ファイル名の長さでソート
	sort.Slice(files, func(i, j int) bool {
		return len(files[i]) < len(files[j])
	})

	return FileGroup{
		keep:     files[0],
		removals: files[1:],
	}
}

func main() {
	// コマンドライン引数からCSVファイルパスを取得
	csvPath := flag.String("csv", "", "重複チェックするCSVファイルのパス")
	outputPath := flag.String("out", "duplicates.txt", "出力テキストファイルのパス")
	debugMode := flag.Bool("debug", false, "デバッグモード（詳細な出力を行う）")
	flag.Parse()

	if *csvPath == "" {
		log.Fatal("CSVファイルのパスを指定してください")
	}

	// CSVファイルをオープン
	file, err := os.Open(*csvPath)
	if err != nil {
		log.Fatalf("CSVファイルのオープンに失敗しました: %v", err)
	}
	defer file.Close()

	csvDir := filepath.Dir(*csvPath)
	outputFullPath := filepath.Join(csvDir, *outputPath)

	r := csv.NewReader(file)
	// CSVファイルはヘッダーがない前提
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("CSVの読み込みに失敗しました: %v", err)
	}

	// SHA256ハッシュをキー、対応するファイル名のスライスを値とするマップを作成
	hashMap := make(map[string][]string)
	for _, rec := range records {
		if len(rec) < 2 {
			continue
		}
		filename := rec[0]
		hash := rec[1]
		hashMap[hash] = append(hashMap[hash], filename)
	}

	// 結果を出力ファイルに書き込み
	outFile, err := os.Create(outputFullPath)
	if err != nil {
		log.Fatalf("出力ファイルの作成に失敗しました: %v", err)
	}
	defer outFile.Close()

	// 重複ファイルグループを処理
	for _, files := range hashMap {
		if len(files) > 1 {
			group := createFileGroup(files)

			if *debugMode {
				// デバッグモード：詳細な出力
				fmt.Fprintf(outFile, "保持するファイル: %s\n", group.keep)
				fmt.Fprintln(outFile, "削除候補:")
				for _, f := range group.removals {
					fmt.Fprintf(outFile, "  %s\n", f)
				}
				fmt.Fprintln(outFile, "---")
			} else {
				// 通常モード：削除候補のみを出力
				for _, f := range group.removals {
					fmt.Fprintln(outFile, f)
				}
			}
		}
	}

	fmt.Printf("重複ファイルのリストを %s に保存しました\n", outputFullPath)
}
