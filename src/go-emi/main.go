package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input loan amount (VND): ")
	principalAmountStr, _ := reader.ReadString('\n')
	fmt.Print("Input interest rate (365 days): ")
	interestRateStr, _ := reader.ReadString('\n')
	fmt.Print("Input loan tenure (months): ")
	loanTenureStr, _ := reader.ReadString('\n')
	fmt.Print("Input first effective date (yyyymmdd): ")
	firstEffectiveDateStr, _ := reader.ReadString('\n')
	fmt.Print("Input first due date (yyyymmdd): ")
	firstDueDateStr, _ := reader.ReadString('\n')
	fmt.Println("----------------------------------------------------------------------------------------------------")
	// principalAmountStr := "30000000"
	// interestRateStr := "26"
	// loanTenureStr := "12"
	// firstEffectiveDateStr := "20210922"
	// firstDueDateStr := "20211016"
	principalAmount, _ := strconv.ParseFloat(strings.ReplaceAll(principalAmountStr, "\n", ""), 64)
	interestRate, _ := strconv.ParseFloat(strings.ReplaceAll(interestRateStr, "\n", ""), 64)
	loanTenure, _ := strconv.Atoi(strings.ReplaceAll(loanTenureStr, "\n", ""))
	firstEffectiveDate, _ := time.Parse("20060102", strings.ReplaceAll(firstEffectiveDateStr, "\n", ""))
	firstDueDate, _ := time.Parse("20060102", strings.ReplaceAll(firstDueDateStr, "\n", ""))
	var effectiveDate time.Time
	var dueDate time.Time
	var days float64
	var portionInterestRate float64
	var pi float64 = 1.0
	var sigma float64
	var portionAmount float64
	for i := 0; i < loanTenure; i++ {
		if i == 0 {
			effectiveDate = firstEffectiveDate
			dueDate = firstDueDate
		} else {
			effectiveDate = firstDueDate.AddDate(0, i-1, 0)
			dueDate = firstDueDate.AddDate(0, i, 0)
		}
		days = dueDate.Sub(effectiveDate).Hours() / 24
		portionInterestRate = interestRate * days / 36500
		pi = pi * 1 / (1 + portionInterestRate)
		sigma = sigma + pi
	}
	fmt.Println(fmt.Sprintf("Loan Amount: %.2f", principalAmount))
	fmt.Println(fmt.Sprintf("Interest Rate: %.2f", interestRate))
	portionAmount = math.Round(principalAmount / sigma)
	fmt.Println(fmt.Sprintf("Portion Amount: %.2f", portionAmount))
	fmt.Println()
	var portionInterestAmount float64
	var remainingPrincipalAmount = principalAmount
	fmt.Println("No.\tEffective Date\tDue Date\tPortion Interest\tPortion Loan Amount\tRemaining Loan Amount")
	for i := 0; i < loanTenure; i++ {
		if i == 0 {
			effectiveDate = firstEffectiveDate
			dueDate = firstDueDate
		} else {
			effectiveDate = firstDueDate.AddDate(0, i-1, 0)
			dueDate = firstDueDate.AddDate(0, i, 0)
		}
		days = dueDate.Sub(effectiveDate).Hours() / 24
		portionInterestRate = interestRate * days / 36500
		portionInterestAmount = math.Round(remainingPrincipalAmount * portionInterestRate)
		fmt.Println(fmt.Sprintf("%d\t%s\t%s\t%.2f\t\t%.2f\t\t%.2f", i+1, effectiveDate.Format("2006-01-02"), dueDate.Format("2006-01-02"), portionInterestAmount, math.Min(portionAmount-portionInterestAmount, remainingPrincipalAmount), remainingPrincipalAmount))
		remainingPrincipalAmount = remainingPrincipalAmount - portionAmount + portionInterestAmount
	}
}
