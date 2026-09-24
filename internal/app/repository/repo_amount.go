package repository

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/yael-castro/budget-tracker/internal/app/business"
	"github.com/yael-castro/budget-tracker/pkg/commavalue"
)

func NewAmountRepositoryV2() business.AmountRepositoryV2 {
	return repository{}
}

type repository struct {
}

func (r repository) GetBudget(ctx context.Context) (business.Budget, error) {
	const (
		fAmount   = 2
		fCategory = 1
	)

	budget := make(business.Budget)

	lineParser := func(record []string) (err error) {
		category := commavalue.Trim(record[fCategory])

		amount, err := commavalue.Float64(record[fAmount])
		if err != nil {
			return
		}

		budget[category] += amount
		return
	}

	err := r.readFile(ctx, budgetFile, lineParser)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

func (r repository) GetExpense(ctx context.Context) (expense business.Expense, err error) {
	const (
		_ = iota - 1
		fDate
		fDescription
		fAmount
		fCategory
	)

	record := business.ExpenseRecord{}

	const expectedRecords = 100
	expense = make(business.Expense, 0, expectedRecords)

	lineParser := func(line []string) (err error) {
		record.Date, err = commavalue.Date(line[fDate])
		if err != nil {
			return
		}

		record.Amount, err = commavalue.Float64(line[fAmount])
		if err != nil {
			return
		}

		record.Category = commavalue.Trim(line[fCategory])
		record.Description = commavalue.Trim(line[fDescription])

		expense = append(expense, record)
		return
	}

	err = r.readFile(ctx, expensesFile, lineParser)
	if err != nil {
		return
	}

	return
}

func (r repository) PutExpense(ctx context.Context, e business.Expense) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) readFile(ctx context.Context, fileName string, f func([]string) error) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	line := []string(nil)
	reader := csv.NewReader(file)

	counter := 0

	for {
		counter++

		line, err = reader.Read()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		// Skipping headers (first line)
		if counter == 1 {
			continue
		}

		err = f(line)
		if err != nil {
			return fmt.Errorf("failed to parse line %d at file %s: %w", counter, fileName, err)
		}
	}
}
