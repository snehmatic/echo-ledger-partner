package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

func main() {
	fmt.Println("Echo ledger partner is here to help!")

	fmt.Println("Reading csv...")

	records := ReadCsv()

	// filter important data
	filter(records)
}

func ReadCsv() []*Record {
	csvFile, csvFileError := os.OpenFile("axio_expense_report_01-01-2020_to_30-11-2025.csv", os.O_RDWR, os.ModePerm)
	if csvFileError != nil {
		panic(csvFileError)
	}
	defer csvFile.Close()

	var articles []*Record
	if unmarshalError := gocsv.UnmarshalFile(csvFile, &articles); unmarshalError != nil {
		panic(unmarshalError)
	}

	return articles
}

func filter(records []*Record) {
	// Category from tracker
	var expAmount float64
	var SpendzSelfShopping float64
	var SpendzFamilyShopping float64
	var SpendzSelfTravel float64
	var SpendzFamilyTravel float64

	year := "2025"

	for _, record := range records {

		val1, err := record.CalculateAmmountByFilter(Filter{Expense: true, Year: year})
		if err != nil {
			panic(err)
		}
		val2, err := record.CalculateAmmountByFilter(Filter{Expense: true, Year: year, Category: "SHOPPING", Tags: []string{"Self"}})
		if err != nil {
			panic(err)
		}
		val3, err := record.CalculateAmmountByFilter(Filter{Expense: true, Year: year, Category: "SHOPPING", Tags: []string{"Family"}})
		if err != nil {
			panic(err)
		}
		val4, err := record.CalculateAmmountByFilter(Filter{Expense: true, Year: year, Category: "TRAVEL", Tags: []string{"Delhi"}})
		if err != nil {
			panic(err)
		}
		val5, err := record.CalculateAmmountByFilter(Filter{Expense: true, Year: year, Category: "TRAVEL", Tags: []string{"Family"}})
		if err != nil {
			panic(err)
		}

		expAmount += val1
		SpendzSelfShopping += val2
		SpendzFamilyShopping += val3
		SpendzSelfTravel += val4
		SpendzFamilyTravel += val5
	}
	// presentation
	fmt.Println("Total expense:", expAmount)

	fmt.Println("Self:")
	fmt.Println("\tShopping:", SpendzSelfShopping)
	fmt.Println("\tTravel:", SpendzSelfTravel)

	fmt.Println("Family:")
	fmt.Println("\tShopping:", SpendzFamilyShopping)
	fmt.Println("\tTravel:", SpendzFamilyTravel)
}
