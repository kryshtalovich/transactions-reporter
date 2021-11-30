package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func createReport(format string) (ReportWriter, error) {
	switch format {
	case "csv":
		return new(csvWriter), nil
	case "json":
		return new(jsonWriter), nil
	default:
		return nil, errors.New("Недопустимый формат файла для сохранения отчета: " + format)
	}
}

func (cw csvWriter) Write(report Report, outputPath string) error {
	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)

	columnsMap := make(map[string]bool)

	for key := range report {
		for field := range report[key] {
			columnsMap[field] = true
		}
	}

	columns := make([]string, 0, len(columnsMap))

	for column := range columnsMap {
		columns = append(columns, column)
	}

	writer.Write(columns)

	for key := range report {
		values := make([]string, 0, len(report[key]))

		for _, field := range columns {
			if val, ok := report[key][field]; ok {
				values = append(values, fmt.Sprint(val))
			}
		}

		writer.Write(values)
	}

	writer.Flush()

	return os.WriteFile(outputPath, []byte(b.String()), 0777)
}

func (jw jsonWriter) Write(report Report, outputPath string) error {
	transactions := make([]UserReport, 0, len(report))

	for _, transaction := range report {
		transactions = append(transactions, transaction)
	}

	data, err := json.MarshalIndent(transactions, "", "\t")

	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0777)
}
