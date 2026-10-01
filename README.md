# Budget tracker
This is a side project just for fun.

I need a way to track my financial habits to achive my financial goals.
## Getting started
### How to install
```shell
go install github.com/yael-castro/budget-tracker/cmd/budget-tracker@latest
```
### How to use
1. Create the file `expenses.csv` 
```text
DATE,DESCRIPTION,AMOUNT,CATEGORY
2026-10-1,MY MONTHLY GROCERIES,200,GROCERIES
```
2. Create the file `budget.csv`
```text
DESCRIPTION,CATEGORY,AMOUNT
Gas,Gas,2000
```
3. Execute the `buget-tracker`
```shell
buget-tracker
```
