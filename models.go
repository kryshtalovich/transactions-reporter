package main

type Report map[int]UserReport

type UserReport map[string]int

type csvWriter struct {
}

type jsonWriter struct {
}

type ReportWriter interface {
	Write(report Report, outputPath string) error
}

type Transaction struct {
	UserId   int    `json:"user_id"`
	Amount   int    `json:"amount"`
	Category string `json:"category"`
}
