package business

import "context"

type (
	ExpensesLogic interface {
		RemainingBudget(context.Context) (AmountTable, error)
	}
)

type (
	AmountRepository interface {
		GetBudget(ctx context.Context) (AmountTable, error)
		GetExpenses(ctx context.Context) (AmountTable, error)
	}
)
