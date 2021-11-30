package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const defaultOutput = "/dev/stdout"
const filename = "./transactions.json"
const reportsDir = "reports/"

var formatFlag = flag.String("f", "json", "Доступные форматы сохранения: csv, json ")
var startTime = time.Now()

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getPath(format string) (string, error) {
	var outFilename string

	if format == "json" || format == "csv" {
		outFilename = reportsDir + "report_" + startTime.Format("01-02-2006_15:04:05") + "." + format
	} else {
		return "", errors.New("Недопустимый формат файла для сохранения отчета: " + format)
	}

	if filepath.IsAbs(outFilename) {
		return outFilename, nil
	}

	dir, err := os.Getwd()
	checkError(err)

	return filepath.Join(dir, outFilename), nil
}

func getPointInfo(text string) {
	fmt.Println(text + ": " + time.Since(startTime).String())
}

func getReportData(file *os.File) (Report, error) {
	value, err := ioutil.ReadAll(file)
	checkError(err)

	getPointInfo("Файл json прочтён")

	var transactions []Transaction
	err = json.Unmarshal(value, &transactions)
	checkError(err)

	getPointInfo("Файл json распаршен")

	report := make(Report)
	for i := 0; i < len(transactions); i++ {
		transaction := transactions[i]

		if report[transaction.UserId] == nil {
			report[transaction.UserId] = make(UserReport)
		}

		report[transaction.UserId]["sum"] += transaction.Amount
		report[transaction.UserId][transaction.Category] += transaction.Amount
		report[transaction.UserId]["user_id"] = transaction.UserId
	}

	return report, nil
}

func main() {
	flag.Parse()

	getPointInfo("Программа запустилась")

	sourceFile, err := os.Open(filename)
	checkError(err)

	defer sourceFile.Close()
	getPointInfo("Файл прочтён")

	report, err := getReportData(sourceFile)
	checkError(err)
	getPointInfo("Отчёт сгенерирован")

	writer, err := createReport(*formatFlag)
	checkError(err)

	outputPath, err := getPath(*formatFlag)
	checkError(err)

	checkError(writer.Write(report, defaultOutput))
	getPointInfo("Отчёт выведен")
	checkError(writer.Write(report, outputPath))
	getPointInfo("Отчёт сохранён")

	getPointInfo("Программа выполнилась")
}
