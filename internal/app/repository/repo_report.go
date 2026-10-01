package repository

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/yael-castro/budget-tracker/internal/app/business"
	"github.com/yael-castro/budget-tracker/pkg/commavalue"
)

const (
	reportFile = "report.csv"
	healthFile = "health.csv"
)

func NewReportRepository() business.FileRepository {
	return reportRepository{}
}

type reportRepository struct{}

func (r reportRepository) SaveFinancialHealth(ctx context.Context, health business.FinancialHealth) error {
	isWrote := false

	saveHealth := func() ([]string, error) {
		if isWrote {
			return nil, io.EOF
		}

		isWrote = true

		return []string{
			commavalue.ToFloat64(health.TotalBudget),
			commavalue.ToFloat64(health.TotalExpenses),
			commavalue.ToFloat64(health.NoBudget),
			commavalue.ToFloat64(health.RemainingBudget),
		}, nil
	}

	headers := []string{"TOTAL BUDGET", "TOTAL EXPENSES", "NO BUDGET", "REMAINING BUDGET"}

	return r.saveFile(ctx, healthFile, headers, saveHealth)
}

func (r reportRepository) BackupExpense(ctx context.Context) error {
	err := os.Rename(expensesFile, expensesFile+".bak")
	if err != nil {
		return err
	}

	return nil
}

func (r reportRepository) SaveExpenseReport(ctx context.Context, report business.ExpenseReport) error {
	index := 0
	line := business.ExpenseReportLine{}

	writeLine := func() ([]string, error) {
		if index >= len(report.Content) {
			return nil, io.EOF
		}

		line = report.Content[index]
		index++

		hasBudget := "Y"
		if line.NoBudget {
			hasBudget = "N"
		}

		return []string{
			line.Category,
			commavalue.ToFloat64(line.Budget),
			commavalue.ToFloat64(line.Expense),
			commavalue.ToFloat64(line.Remaining),
			hasBudget,
		}, nil
	}

	headers := []string{"Category", "Budget", "Expense", "Remaining", "Has Budget"}

	return r.saveFile(ctx, reportFile, headers, writeLine)
}

func (r reportRepository) SaveExpense(ctx context.Context, expense business.Expense) error {
	index := 0
	line := business.ExpenseRecord{}

	writeLine := func() ([]string, error) {
		if index >= len(expense) {
			return nil, io.EOF
		}

		line = expense[index]
		index++

		return []string{
			line.Date.Format("2006-01-02"), // TODO: add function at package level commavalue
			commavalue.Trim(line.Description),
			commavalue.ToFloat64(line.Amount),
			commavalue.Trim(line.Category),
		}, nil
	}

	headers := []string{"DATE", "DESCRIPTION", "AMOUNT", "CATEGORY"}
	return r.saveFile(ctx, expensesFile, headers, writeLine)
}

func (r reportRepository) saveFile(ctx context.Context, fileName string, headers []string, f func() ([]string, error)) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	writer := csv.NewWriter(file)

	err = writer.Write(headers)
	if err != nil {
		return err
	}
	defer writer.Flush()

	var record []string

	for {
		record, err = f()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}
}
