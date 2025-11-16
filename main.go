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

	var records []*Record
	if unmarshalError := gocsv.UnmarshalFile(csvFile, &records); unmarshalError != nil {
		panic(unmarshalError)
	}

	return records
}

func filter(records []*Record) {
	// Category from tracker
	var expenseTotal float64
	var incomeTotal float64
	var selfShopping float64
	var familyShopping float64
	var selfTravel float64
	var familyTravel float64

	year := "2025"

	for _, record := range records {

		expense, err := record.CalculateAmmountByFilter(Filter{Expense: true, Year: year})
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

		income, err := record.CalculateAmmountByFilter(Filter{Income: true, Year: year})
		if err != nil {
			panic(err)
		}

		expenseTotal += expense
		incomeTotal += income
		selfShopping += val2
		familyShopping += val3
		selfTravel += val4
		familyTravel += val5
	}

	// presentation
	fmt.Printf("Total Income: %.2f\n", incomeTotal)
	fmt.Printf("Total Expense: %.2f\n", expenseTotal)
	fmt.Printf("Difference: %.2f\n", incomeTotal-expenseTotal)

	fmt.Println("Self:")
	fmt.Printf("\tShopping: %.2f\n", selfShopping)
	fmt.Printf("\tTravel: %.2f\n", selfTravel)

	fmt.Println("Family:")
	fmt.Printf("\tShopping: %.2f\n", familyShopping)
	fmt.Printf("\tTravel: %.2f\n", familyTravel)
}
