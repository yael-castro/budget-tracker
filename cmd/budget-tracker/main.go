package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yael-castro/budget-tracker/internal/app/business"
	"github.com/yael-castro/budget-tracker/internal/app/repository"
)

func main() {
	// Building main context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Building dependencies
	reportRepo := repository.NewReportRepository()
	amountRepo := repository.NewAmountRepositoryV2()
	reportLogic := business.NewReportLogic(reportRepo, amountRepo)

	// Executing command
	err := reportLogic.ReportExpense(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
