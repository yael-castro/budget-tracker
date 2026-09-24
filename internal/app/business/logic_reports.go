package business

import "context"

func NewReportLogic(fileRepo FileRepository, amountRepo AmountRepositoryV2) reportLogic {
	return reportLogic{
		fileRepo:   fileRepo,
		amountRepo: amountRepo,
	}
}

type reportLogic struct {
	fileRepo   FileRepository
	amountRepo AmountRepositoryV2
}

func (l reportLogic) ReportExpense(ctx context.Context) error {
	// Getting financial information
	budget, err := l.amountRepo.GetBudget(ctx)
	if err != nil {
		return err
	}

	expense, err := l.amountRepo.GetExpense(ctx)
	if err != nil {
		return err
	}

	// Generating expense report
	expenseReport := NewExpenseReport(budget, expense)

	// Generating backup of expense file
	err = l.fileRepo.BackupExpense(ctx)
	if err != nil {
		return err
	}

	// Saving expense report
	err = l.fileRepo.SaveExpenseReport(ctx, expenseReport)
	if err != nil {
		return err
	}

	// Updating expense file
	err = l.fileRepo.SaveExpense(ctx, expense.Order())
	if err != nil {
		return err
	}

	// Saving financial health
	health := expenseReport.FinancialHealth()

	err = l.fileRepo.SaveFinancialHealth(ctx, health)
	if err != nil {
		return err
	}

	return nil
}
