package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/yael-castro/budget-tracker/internal/app/business"
	"github.com/yael-castro/budget-tracker/internal/app/driving/csvfile"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// Setting logger
	log.Default().SetFlags(0)

	// Building dependencies
	amountRepo := csvfile.NewAmountRepository()
	expensesLogic := business.NewExpensesLogic(amountRepo)

	// Executing logic
	remaining, err := expensesLogic.RemainingBudget(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(remaining)
}
