package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

// file拡張子を削除する
func getFileNameNoExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func main() {
	var csvfile string
	var title bool
	var lines int

	flag.StringVar(&csvfile, "f", "", "CSV file name")
	flag.BoolVar(&title, "t", false, "1行目をタイトル行とする。")
	flag.IntVar(&lines, "l", 1, "分割するライン数を指定します。")
	flag.Parse()

	inFile, err := os.Open(csvfile)
	failOnError(err)
	defer inFile.Close()

	reader := csv.NewReader(transform.NewReader(inFile, japanese.ShiftJIS.NewDecoder()))
	reader.LazyQuotes = true

	log.Printf("Start CSV Reading...")

	var titleLine []string

	if title {
		titleLine, err = reader.Read()
		failOnError(err)
	}

	var outFile *os.File
	var fileopenerr error
	var writer *csv.Writer
	var csvOutFile string

	for {
		for i := 0; i < lines; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				println("end")
				os.Exit(1)
			} else {
				failOnError(err)
			}

			if i == 0 {

				// ディレクトリ作成
				os.Mkdir("tmp", 0777)

				//fmt.Println(record)
				csvOutFile = getFileNameNoExt(csvfile)
				if lines == 1 {
					csvOutFile = "tmp/" + csvOutFile + "_" + record[1] + ".csv"
				} else {
					csvOutFile = "tmp/" + csvOutFile + "_" + record[1] + "_" + strconv.Itoa(lines) + "_in_records" + ".csv"
				}

				outFile, fileopenerr = os.Create(csvOutFile)
				failOnError(fileopenerr)

			}
			//write csv in utf-8
			writer = csv.NewWriter(outFile)
			if title && i == 0 {
				writer.Write(titleLine)
			}
			writer.Write(record)
			writer.Flush()
		}
		outFile.Close()
	}
}
