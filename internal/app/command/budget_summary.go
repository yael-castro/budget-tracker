package command

import (
	"context"
	"io"
	"os"

	"github.com/yael-castro/budget-tracker/internal/app/business"
)

func BudgetSummary(logic business.ExpensesLogic) func(context.Context) error {
	return func(ctx context.Context) (err error) {
		summary, err := logic.GenerateBudgetSummary(ctx)
		if err != nil {
			return
		}

		_, err = io.Copy(os.Stdout, summary.File())
		return
	}
}
