package business

import "context"

type (
	ExpensesLogic interface {
		GenerateBudgetSummary(context.Context) (BudgetSummary, error)
	}
)

type (
	AmountRepository interface {
		GetBudget(ctx context.Context) (AmountTable, error)
		GetExpenses(ctx context.Context) (AmountTable, error)
	}

	AmountRepositoryV2 interface {
		GetBudget(context.Context) (Budget, error)
		GetExpense(context.Context) (Expense, error)
	}

	FileRepository interface {
		BackupExpense(ctx context.Context) error
		SaveExpense(context.Context, Expense) error
		SaveExpenseReport(context.Context, ExpenseReport) error
		SaveFinancialHealth(context.Context, FinancialHealth) error
	}
)
