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

func (e expensesLogic) RemainingBudget(ctx context.Context) (remaining AmountTable, err error) {
	var manager Manager

	manager.Budget, err = e.amountRepo.GetBudget(ctx)
	if err != nil {
		return
	}

	manager.Expenses, err = e.amountRepo.GetExpenses(ctx)
	if err != nil {
		return
	}

	remaining = manager.RemainingBudget()
	return
}
