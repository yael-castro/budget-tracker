package business

import (
	"bytes"
	"encoding/csv"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Types for avoiding have naked params
type (
	Category = string
	Amount   = float64
)

type BudgetSummary struct {
	Summary
}

func (s BudgetSummary) File() io.ReadCloser {
	// Creating output
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Writing headers
	_ = writer.Write([]string{"CATEGORY", "BUDGET", "AMOUNT", "REMAINING", "IS_CONTEMPLATED", "IS_EXCEED"})

	// Calculating remaining budget for each category
	budget := s.RemainingBudget()

	totalBudget := 0.0
	totalExpense := 0.0
	lackBudget := 0.0

	// Translating records from
	for category, remaining := range budget {
		budget := s.Budget[category]
		totalBudget += budget

		expense := s.Expenses[category]
		totalExpense += expense

		if budget <= 0 {
			lackBudget += expense
		}

		// Writing a remaining budget's record for a category
		_ = writer.Write([]string{
			category,
			strconv.FormatFloat(budget, 'f', -1, 64),
			strconv.FormatFloat(expense, 'f', -1, 64),
			strconv.FormatFloat(remaining, 'f', -1, 64),
			strconv.FormatBool(budget > 0),
			strconv.FormatBool(expense > budget),
		})
	}

	// Writing a record to show total budget
	_ = writer.Write([]string{
		"TOTAL BUDGET",
		strconv.FormatFloat(totalBudget, 'f', -1, 64),
		strconv.FormatFloat(totalExpense, 'f', -1, 64),
		strconv.FormatFloat(totalBudget-totalExpense, 'f', -1, 64),
	})

	// Writing a record to show lack budget
	_ = writer.Write([]string{
		"LACK BUDGET",
		strconv.FormatFloat(lackBudget, 'f', -1, 64),
	})

	writer.Flush()
	return io.NopCloser(buf)
}

func (s BudgetSummary) String() string {
	var b strings.Builder

	// Ordering keys
	budget := s.Summary.RemainingBudget()

	b.WriteString("Category,Budget,Used,Remaining\n")

	totalBudget := 0.0
	totalExpense := 0.0
	lackBudget := 0.0

	for category, remaining := range budget {
		budget := s.Budget[category]
		totalBudget += budget

		expense := s.Expenses[category]
		totalExpense += expense

		if budget <= 0 {
			lackBudget += expense
		}

		b.WriteString(category)
		b.WriteRune(',')
		b.WriteString(strconv.FormatFloat(budget, 'f', -1, 64))
		b.WriteRune(',')
		b.WriteString(strconv.FormatFloat(expense, 'f', -1, 64))
		b.WriteRune(',')
		b.WriteString(strconv.FormatFloat(remaining, 'f', -1, 64))
		b.WriteRune('\n')
	}

	b.WriteString("TOTAL BUDGET")
	b.WriteRune(',')
	b.WriteString(strconv.FormatFloat(totalBudget, 'f', -1, 64))
	b.WriteRune(',')
	b.WriteString(strconv.FormatFloat(totalExpense, 'f', -1, 64))
	b.WriteRune(',')
	b.WriteString(strconv.FormatFloat(totalBudget-totalExpense, 'f', -1, 64))
	b.WriteRune('\n')

	//
	b.WriteString("LACK BUDGET")
	b.WriteRune(',')
	b.WriteString(strconv.FormatFloat(lackBudget, 'f', -1, 64))
	b.WriteRune(',')
	b.WriteString("...")
	b.WriteRune(',')
	b.WriteString("...")
	b.WriteRune('\n')

	return b.String()
}

type Summary struct {
	Budget   AmountTable
	Expenses AmountTable
}

func (m Summary) RemainingBudget() AmountTable {
	return m.Budget.Minus(m.Expenses)
}

type AmountRecord struct {
	Day        int
	Month      int
	Year       int
	IsRequired bool
	Amount     float64
}

type AmountTable map[Category]Amount

func (a AmountTable) Minus(b AmountTable) AmountTable {
	keys := make(map[Category]struct{}, len(a))

	for k := range a {
		keys[k] = struct{}{}
	}

	for k := range b {
		keys[k] = struct{}{}
	}

	// Resting
	c := make(map[Category]Amount, len(keys))

	for k := range keys {
		c[k] = a[k] - b[k]
	}

	return c
}

func (a AmountTable) Total() (total Amount) {
	for _, amount := range a {
		total += amount
	}

	return
}

func (a AmountTable) Amount(category Category) Amount {
	return a[category]
}

func (a AmountTable) String() string {
	var b strings.Builder

	b.WriteString("Category,Amount\n")
	for category, amount := range a {
		b.WriteString(category)
		b.WriteString(",")
		b.WriteString(strconv.FormatFloat(amount, 'f', 2, 64))
		b.WriteRune('\n')
	}

	return b.String()
}

type Budget = AmountTable
type Expense []ExpenseRecord

func (e Expense) Table() AmountTable {
	table := make(AmountTable, len(e))

	for _, expense := range e {
		table[expense.Category] += expense.Amount
	}

	return table
}

func (e Expense) Order() Expense {
	sort.Slice(e, func(left, right int) bool {
		if e[left].Date.Equal(e[right].Date) {
			return e[left].Amount < e[right].Amount
		}

		return e[left].Date.Before(e[right].Date)
	})

	return e
}

type ExpenseRecord struct {
	Date        time.Time
	Amount      Amount
	Category    Category
	Description string
}

func NewExpenseReport(budget Budget, expense Expense) ExpenseReport {
	wasted := expense.Table()
	joined := budget.Minus(wasted)

	report := ExpenseReport{
		Content: make([]ExpenseReportLine, 0, len(joined)),
	}

	for category, amount := range joined {
		_, budgetExists := budget[category]

		line := ExpenseReportLine{
			Category:  category,
			Budget:    budget[category],
			Expense:   wasted[category],
			Remaining: amount,
			NoBudget:  !budgetExists,
		}

		report.Content = append(report.Content, line)
	}

	return report
}

type FinancialHealth struct {
	NoBudget        Amount
	TotalBudget     Amount
	TotalExpenses   Amount
	RemainingBudget Amount
}

type ExpenseReport struct {
	Content []ExpenseReportLine
}

func (r ExpenseReport) FinancialHealth() FinancialHealth {
	health := FinancialHealth{}

	for _, line := range r.Content {
		health.TotalBudget += line.Budget
		health.TotalExpenses += line.Expense

		if line.NoBudget {
			health.NoBudget += line.Expense
		}
	}

	health.RemainingBudget = health.TotalBudget - health.TotalExpenses

	return health
}

type ExpenseReportLine struct {
	Category  Category
	Budget    Amount
	Expense   Amount
	Remaining Amount
	NoBudget  bool
}
