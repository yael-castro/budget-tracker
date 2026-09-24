package business

import "context"

func NewExpensesLogic(amountRepo AmountRepository) ExpensesLogic {
	return expensesLogic{
		amountRepo: amountRepo,
	}
}

type expensesLogic struct {
	amountRepo AmountRepository
}

func (e expensesLogic) GenerateBudgetSummary(ctx context.Context) (remaining BudgetSummary, err error) {
	remaining.Budget, err = e.amountRepo.GetBudget(ctx)
	if err != nil {
		return
	}

	remaining.Expenses, err = e.amountRepo.GetExpenses(ctx)
	if err != nil {
		return
	}

	return
}
