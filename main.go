package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	year      = "2025"
	dateRange = DateRange{
		Start: "01-01-2025",
		End:   "31-12-2025",
	}

	CategoryWiseFilterMap = map[string]Filter{
		"Income Total":        {Income: true, DateRange: &dateRange, Year: year},
		"Expense Total":       {Expense: true, DateRange: &dateRange, Year: year},
		"Delhi (Self Travel)": {Expense: true, DateRange: &dateRange, Year: year, Tags: []string{"Delhi"}},
		"Self (Shopping)":     {Expense: true, DateRange: &dateRange, Year: year, Category: "SHOPPING", Tags: []string{"Self"}},
		"Family (Travel)":     {Expense: true, DateRange: &dateRange, Year: year, Category: "TRAVEL", Tags: []string{"Family"}},
		"Family (Shopping)":   {Expense: true, DateRange: &dateRange, Year: year, Category: "SHOPPING", Tags: []string{"Family"}},
	}

	PresentationMap = make(map[string]float64)
)

func main() {
	fmt.Println("Echo ledger partner is here to help!")

	fmt.Println("Reading csv...")

	records := ReadCsv()

	fmt.Println("------------------------------------")
	// filter and print important data
	filter(records)
	fmt.Println("------------------------------------")
}

func ReadCsv() []*Record {
	csvFile, csvFileError := os.OpenFile("axio_expense_report_01-01-2020_to_30-11-2025.csv", os.O_RDWR, os.ModePerm)
	if csvFileError != nil {
		panic(csvFileError)
	}
	defer csvFile.Close()

	var records []*Record
	if unmarshalError := gocsv.UnmarshalFile(csvFile, &records); unmarshalError != nil {
		panic(unmarshalError)
	}

	return records
}

func filter(records []*Record) {

	for _, record := range records {

		for catName, catFilter := range CategoryWiseFilterMap {
			val, err := record.CalculateAmmountByFilter(catFilter)
			if err != nil {
				panic(err)
			}

			PresentationMap[catName] += val
		}
	}

	// presentation
	p := message.NewPrinter(language.Hindi)

	for k := range CategoryWiseFilterMap {
		p.Printf("%s:  %.2f\n", k, PresentationMap[k])
	}
}
