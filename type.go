package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	layout = "02-01-2006" // Corresponds to DD-MM-YYYY
)

type Record struct {
	DateStr   string `csv:"DATE"`
	TimeStr   string `csv:"TIME"`
	Place     string `csv:"PLACE"`
	AmountStr string `csv:"AMOUNT"`
	Kind      string `csv:"DR/CR"`
	Account   string `csv:"ACCOUNT"`
	Expense   string `csv:"EXPENSE"`
	Income    string `csv:"INCOME"`
	Category  string `csv:"CATEGORY"`
	Tags      string `csv:"TAGS"`
	Note      string `csv:"NOTE"`
}

type Filter struct {
	Expense   bool   // is expense
	Income    bool   // is income
	Category  string // category substring
	Year      string
	DateRange *DateRange
	Tags      []string // tag substrings that need to be in record
}

type DateRange struct {
	Start string // ex. 30-01-2000
	End   string // ex. 30-12-2025
}

func (dr DateRange) ParsedDateRange() (time.Time, time.Time) {
	start, _ := time.Parse(layout, dr.Start)
	end, _ := time.Parse(layout, dr.End)
	return start, end
}

func (r Record) IsInDateRange(dr DateRange) bool {
	rDate, _ := r.ParsedDate()
	start, end := dr.ParsedDateRange()

	if start.After(end) {
		start, end = end, start
	}

	return (rDate.Equal(start) || rDate.After(start)) && (rDate.Equal(end) || rDate.Before(end))
}

func (r Record) Amount() (float64, error) {
	// convert amount string to float64 (ex. 1,000 -> 1000.00)
	cleanedStr := strings.ReplaceAll(r.AmountStr, ",", "")
	cleanedStr = strings.ReplaceAll(cleanedStr, "'", "")
	amount, err := strconv.ParseFloat(cleanedStr, 64)
	if err != nil {
		return 0.00, fmt.Errorf("Error converting string to float: %v\n", err)
	}

	return amount, nil
}

func (r Record) ParsedDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.DateStr)
}

func (r Record) CalculateAmmountByFilter(filter Filter) (float64, error) {

	amount, err := r.Amount()
	if err != nil {
		return 0.00, err
	}

	if filter.Expense {
		if !strings.Contains(r.Expense, "Yes") {
			return 0, nil
		}
	} else if filter.Income {
		if !strings.Contains(r.Income, "Yes") {
			return 0, nil
		}
	}

	if filter.Year != "" {
		if !strings.Contains(r.DateStr, filter.Year) {
			return 0, nil
		}
	}

	if filter.DateRange != nil {
		if !r.IsInDateRange(*filter.DateRange) {
			return 0, nil
		}
	}

	if filter.Category != "" && r.Category != filter.Category {
		return 0, nil
	}

	if len(filter.Tags) > 0 {
		for _, tag := range filter.Tags {
			if !strings.Contains(r.Tags, tag) {
				// if even one of the filter tags aren't present, return
				return 0, nil
			}
		}
	}

	return amount, nil
}
