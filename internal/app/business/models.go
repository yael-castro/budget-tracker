package business

import (
	"strconv"
	"strings"
)

// Types for avoiding have naked params
type (
	Category = string
	Amount   = float64
)

type Manager struct {
	Budget   AmountTable
	Expenses AmountTable
}

func (m Manager) RemainingBudget() AmountTable {
	return m.Budget.Minus(m.Expenses)
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
