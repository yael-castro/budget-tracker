package csvfile

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/yael-castro/budget-tracker/internal/app/business"
)

const (
	budgetFile   = "budget.csv"
	expensesFile = "expenses.csv"
)

func NewAmountRepository() business.AmountRepository {
	return amountRepository{}
}

type amountRepository struct{}

func (a amountRepository) GetBudget(ctx context.Context) (business.AmountTable, error) {
	const categoryColumn, amountColumn = 1, 2

	table, err := a.readFile(ctx, budgetFile, categoryColumn, amountColumn)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget table: %w", err)
	}

	return table, nil
}

func (a amountRepository) GetExpenses(ctx context.Context) (business.AmountTable, error) {
	const categoryColumn, amountColumn = 3, 2

	table, err := a.readFile(ctx, expensesFile, categoryColumn, amountColumn)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses table: %w", err)
	}

	return table, nil
}

func (a amountRepository) readFile(_ context.Context, fileName string, category, amount int) (table business.AmountTable, err error) {
	// Opening file
	budgetFile, err := os.Open(fileName)
	if err != nil {
		return
	}

	// Closing file at the end
	defer func() {
		_ = budgetFile.Close()
	}()

	// Reading file
	reader := csv.NewReader(budgetFile)

	// In memory func
	upper := func(str string) string {
		return strings.ToUpper(strings.TrimSpace(str))
	}

	// Reading all rows
	rows, err := reader.ReadAll()
	if err != nil {
		return
	}

	// Parsing columns and rows in a business.AmountTable
	table = make(business.AmountTable)

	if len(rows) <= 1 {
		return
	}

	for _, row := range rows[1:] {
		categoryColumn := upper(row[category])
		amountColumn := upper(row[amount])

		table[categoryColumn], err = strconv.ParseFloat(amountColumn, 64)
		if err != nil {
			return
		}
	}

	return
}
