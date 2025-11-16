package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	csvPath := flag.String("path", "", "path to csv file")

	if *csvPath == "" {
		fmt.Println("Path to csv file was not provided, getting csv files in current directory...")
		files := findCSVFiles(".")
		fmt.Println("Using first csv file as default...")
		csvPath = &files[0]

	}

	fmt.Printf("Reading csv file (%s)...\n", *csvPath)

	records := ReadCsv(*csvPath)

	fmt.Println("---------------------------------------------------------")
	// filter and print important data
	filter(records)
	fmt.Println("---------------------------------------------------------")
}

func ReadCsv(csvPath string) []*Record {
	csvFile, csvFileError := os.OpenFile(csvPath, os.O_RDWR, os.ModePerm)
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

func findCSVFiles(dir string) []string {
	var fileNames []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	fmt.Printf("CSV files found in %s:\n", dir)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			fmt.Println("\t", len(fileNames)+1, ".", filepath.Join(dir, file.Name()))
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}
