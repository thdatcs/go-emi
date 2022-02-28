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
	loanAmountStr, _ := reader.ReadString('\n')
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
	loanAmount, _ := strconv.ParseFloat(strings.ReplaceAll(loanAmountStr, "\n", ""), 64)
	interestRate, _ := strconv.ParseFloat(strings.ReplaceAll(interestRateStr, "\n", ""), 64)
	loanTenure, _ := strconv.Atoi(strings.ReplaceAll(loanTenureStr, "\n", ""))
	firstEffectiveDate, _ := time.Parse("20060102", strings.ReplaceAll(firstEffectiveDateStr, "\n", ""))
	firstDueDate, _ := time.Parse("20060102", strings.ReplaceAll(firstDueDateStr, "\n", ""))
	portionAmount := calculatePortionAmount(loanAmount, interestRate, loanTenure, firstEffectiveDate, firstDueDate)
	detailPlans := calculateDetailPlans(loanAmount, interestRate, loanTenure, firstEffectiveDate, firstDueDate, portionAmount)
	fmt.Println(fmt.Sprintf("Loan Amount: %.2f", loanAmount))
	fmt.Println(fmt.Sprintf("Interest Rate: %.2f", interestRate))
	fmt.Println(fmt.Sprintf("Portion Amount: %.2f", portionAmount))
	fmt.Println()
	fmt.Println("No.\tEffective Date\tDue Date\tPortion Interest\tPortion Loan Amount\tRemaining Loan Amount")
	for i := 0; i < len(detailPlans); i++ {
		fmt.Println(fmt.Sprintf("%d\t%s\t%s\t%s\t\t%s\t\t%s", i+1, detailPlans[i]["Effective Date"], detailPlans[i]["Due Date"], detailPlans[i]["Portion Interest Amount"], detailPlans[i]["Portion Loan Amount"], detailPlans[i]["Remaining Loan Amount"]))

	}
}

func calculatePortionAmount(loanAmount, interestRate float64, loanTenure int, firstEffectiveDate, firstDueDate time.Time) float64 {
	var effectiveDate time.Time
	var dueDate time.Time
	var days float64
	var portionInterestRate float64
	var pi float64 = 1.0
	var sigma float64 = 0.0
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
	portionAmount = math.Round(loanAmount / sigma)
	return portionAmount
}

func calculateDetailPlans(loanAmount, interestRate float64, loanTenure int, firstEffectiveDate, firstDueDate time.Time, portionAmount float64) []map[string]string {
	var effectiveDate time.Time
	var dueDate time.Time
	var days float64
	var portionInterestRate float64
	var portionInterestAmount float64
	var prevRemainingLoanAmount float64
	var remainingLoanAmount = loanAmount
	var detailPlans []map[string]string
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
		portionInterestAmount = math.Round(remainingLoanAmount * portionInterestRate)
		prevRemainingLoanAmount = remainingLoanAmount
		remainingLoanAmount = remainingLoanAmount - portionAmount + portionInterestAmount
		var detailPlan = make(map[string]string)
		detailPlan["Effective Date"] = effectiveDate.Format("2006-01-02")
		detailPlan["Due Date"] = dueDate.Format("2006-01-02")
		detailPlan["Portion Interest Amount"] = fmt.Sprintf("%.2f", portionInterestAmount)
		detailPlan["Portion Loan Amount"] = fmt.Sprintf("%.2f", math.Min(portionAmount-portionInterestAmount, prevRemainingLoanAmount))
		detailPlan["Remaining Loan Amount"] = fmt.Sprintf("%.2f", math.Max(0, remainingLoanAmount))
		detailPlans = append(detailPlans, detailPlan)
	}
	return detailPlans
}
